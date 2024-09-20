package src

import "github.com/bwmarrin/discordgo"

func CreateResponseDataSignup() *discordgo.MessageSend {
	btns := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "Game_1",
			Style:    1,
			Disabled: false,
			URL:      "",
			CustomID: "game_1_btn",
		},
		discordgo.Button{
			Label:    "Game_2",
			Style:    1,
			Disabled: false,
			CustomID: "game_2_btn",
		},
		discordgo.Button{
			Label:    "Game_3",
			Style:    1,
			Disabled: false,
			CustomID: "game_3_btn",
		},
	}

	embeds := GetEmbeds()

	return &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: btns,
			},
		},
		Embeds: embeds,
	}
}

func GetEmbeds() []*discordgo.MessageEmbed {
	// TODO add times from .env
	return []*discordgo.MessageEmbed{
		{
			Type:  "rich",
			Title: `Game_1: 19:30`,
			Color: 0xff00ae,
		},
		{
			Type:  "rich",
			Title: `Game_2: 20:45`,
			Color: 0xff00ae,
		},
		{
			Type:  "rich",
			Title: `Game_3: 22:00`,
			Color: 0xff00ae,
		},
	}
}
