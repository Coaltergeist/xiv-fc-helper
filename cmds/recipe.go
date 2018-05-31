package cmds

import (
	"github.com/bwmarrin/discordgo"
)

func recipeCommand(s *discordgo.Session, m *discordgo.Message) {
	s.ChannelMessageSend(m.ChannelID, "Coming Soon!")
}

func init() {
	add(&command{
		execute: recipeCommand,
		trigger: "recipe",
		aliases: []string{"craft"},
		desc:    "See needed mats",
	})
}
