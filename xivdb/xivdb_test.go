package xivdb

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/coaltergeist/xiv-fc-helper/structs"
	"github.com/stretchr/testify/assert"
)

func TestCharacterSearch(t *testing.T) {
	name := "neko mcdoogal"
	request := NewSearchRequest()
	request.SetSearch(name)
	request.SetType(CHARACTER)

	data := request.Queue().Consume()

	var characterResults structs.XIVDBCharacterSearch
	if err := json.Unmarshal(data, &characterResults); err != nil {
		panic(err)
	}
	assert.Equal(t, characterResults.Characters.Total, 1, "# of results found not equal to 1")
}

func TestCharacterQuery(t *testing.T) {
	id := 12285614
	request := NewQueryRequest()
	request.SetID(id)
	request.SetType(CHARACTER)

	data := request.Queue().Consume()
	var character structs.XIVDBCharacter
	if err := json.Unmarshal(data, &character); err != nil {
		panic(err)
	}
	strID := strconv.Itoa(id)
	assert.Equal(t, strID, character.Data.ID, "found ID not equal to queried ID")
}
