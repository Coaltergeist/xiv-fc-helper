package cmds

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/coaltergeist/xiv-fc-helper/xivdb"
)

func recipeCommand(s *discordgo.Session, m *discordgo.Message) {
	split := strings.Split(m.Content, " ")

	query := ""
	for i := range split {
		if i != 0 {
			query += split[i] + " "
		}
	}
	//s.ChannelMessageSend(m.ChannelID, "Looking up "+query+"now!")
	request := xivdb.NewSearchRequest()
	request.SetType(xivdb.RECIPE)
	request.SetSearch(fmt.Sprintf(query))
	data := request.Queue().Consume()
}

func init() {
	add(&command{
		execute: recipeCommand,
		trigger: "recipe",
		aliases: []string{"craft"},
		desc:    "See needed mats",
	})
}
