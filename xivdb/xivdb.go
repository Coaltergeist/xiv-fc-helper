// Package xivdb is a framework for making requests to XIVDB's API
// while following their rate limit limitations
package xivdb

import (
	"io/ioutil"
	"net/http"
)

// A Request is a struct nested that provides an endpoint,
// type, and return data channel
type Request struct {
	requestType RequestType
	endpoint    string

	data chan ([]byte)
}

// A SearchRequest is a request to search for something in XIVDB
type SearchRequest struct {
	request      *Request
	searchString string
}

// A QueryRequest is a request to query a single resource
type QueryRequest struct {
	request *Request
	id      int
}

const (
	// CHARACTER is a character in FFXIV
	CHARACTER = iota
	// ITEM is an item in FFXIV
	ITEM
	// RECIPE is a crafting recipe in FFXIV
	RECIPE
)

// A RequestType is what type of request we're making, i.e. character, item etc
type RequestType int

// Queue queues the request into the request queue
func (r *SearchRequest) Queue() *Request {
	r.generateEndpoint()
	requestQueue.push(r.request)
	return r.request
}

// Queue queues the request into the request queue
func (r *QueryRequest) Queue() *Request {
	r.generateEndpoint()
	requestQueue.push(r.request)
	return r.request
}

// GetEndpoint returns the XIVDB endpoint of this request. This isn't finalized until Queue() is called
func (r *SearchRequest) GetEndpoint() string {
	return r.request.endpoint
}

// GetEndpoint returns the XIVDB endpoint of this request. This isn't finalized until Queue() is called
func (r *QueryRequest) GetEndpoint() string {
	return r.request.endpoint
}

// Execute executes a request and returns the string response of the body
func (r *Request) execute() {
	defer func() {
		if r := recover(); r != nil {
			l.Println("recovered in request Execute")
		}
	}()
	resp, err := http.Get(r.endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	r.data <- body
	close(r.data)
}

// Consume is a blocking operation that attempts to consume the request's response
func (r *Request) Consume() []byte {
	return <-r.data
}

// NewQueryRequest creates a new request to a well formed, specific resource
// (i.e. not a search)
func NewQueryRequest() *QueryRequest {
	return &QueryRequest{
		id: -1,
		request: &Request{
			requestType: -1,
			data:        make(chan []byte),
		},
	}
}

// NewSearchRequest creates a new request to search for a resource
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{
		request: &Request{
			requestType: -1,
			data:        make(chan []byte),
		},
	}
}

// SetType sets the type of search
func (r *SearchRequest) SetType(requestType RequestType) {
	r.request.requestType = requestType
}

// SetType sets the type of search
func (r *QueryRequest) SetType(requestType RequestType) {
	r.request.requestType = requestType
}

// SetSearch sets the search string to look for
func (r *SearchRequest) SetSearch(searchString string) {
	r.searchString = searchString
}

// SetID sets the ID to query for
func (r *QueryRequest) SetID(queryID int) {
	r.id = queryID
}
