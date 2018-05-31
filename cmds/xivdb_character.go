package cmds

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/coaltergeist/xiv-fc-helper/structs"
	"github.com/coaltergeist/xiv-fc-helper/xivdb"
	"github.com/paul-io/discordgo-embeds/colors"
	"github.com/paul-io/discordgo-embeds/embed"
)

func xivdbCharacterSearchCommand(s *discordgo.Session, m *discordgo.Message) {
	split := strings.Split(m.Content, " ")
	if len(split) != 4 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Incorrect usage. Use %shelp xivdbcharacter for correct usage", commandPrefix))
		return
	}

	// world := split[1]
	first := split[2]
	last := split[3]

	request := xivdb.NewSearchRequest()
	request.SetType(xivdb.CHARACTER)
	request.SetSearch(fmt.Sprintf("%s+%s", first, last))
	data := request.Queue().Consume()

	var characterResults structs.XIVDBCharacterSearch
	if err := json.Unmarshal(data, &characterResults); err != nil {
		panic(err)
	}
	if characterResults.Characters.Total == 0 || characterResults.Characters.Total > 1 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Non-1 result count: %d\n<%s>", characterResults.Characters.Total, request.GetEndpoint()))
		return
	}

	charRequest := xivdb.NewQueryRequest()
	charRequest.SetID(characterResults.Characters.Results[0].ID)
	charRequest.SetType(xivdb.CHARACTER)

	data = charRequest.Queue().Consume()
	var character structs.XIVDBCharacter
	if err := json.Unmarshal(data, &character); err != nil {
		panic(err)
	}

	em := embed.New()

	description := make([]string, 0, 3)
	description = append(description, fmt.Sprintf("**%s** | **%s**", character.Data.Race, character.Data.Clan))
	description = append(description, fmt.Sprintf("**Nameday**: %s | **Guardian**: %s", character.Data.Nameday, character.Data.Guardian.Name))
	if character.Data.GrandCompany != nil {
		description = append(description, fmt.Sprintf("**Grand Company**: %s / %s", character.Data.GrandCompany.Name, character.Data.GrandCompany.Rank))
	}
	// description = append(description, fmt.Sprintf("**Free Company**: %s", character.Data.))

	em.SetAuthor(fmt.Sprintf("%s of %s", character.Name, character.Data.Server), character.URLLodestone, character.Data.ActiveClass.Role.Icon).
		SetImage(character.Portrait).
		SetDesc(strings.Join(description, "\n")).
		SetColor(colors.Cyan())

	s.ChannelMessageSendEmbed(m.ChannelID, em.MessageEmbed)
}

func init() {
	add(&command{
		execute: xivdbCharacterSearchCommand,
		trigger: "xivdbcharacter",
		aliases: []string{"xivdbchar"},
		desc:    "Search for a character on xivdb",
		usage:   "xivdbcharacter first-name last-name",
	})
}
