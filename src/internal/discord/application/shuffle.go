package application

import (
	"context"
	"dota-discord-bot/src/internal/kungdota/service"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	ShuffleCmd = discordgo.ApplicationCommand{
		Name:        "shuffle",
		Description: "shuffle teams",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_0",
				Description: "player 0",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_1",
				Description: "player 1",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_2",
				Description: "player 2",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_3",
				Description: "player 3",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_4",
				Description: "player 4",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_5",
				Description: "player 5",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_6",
				Description: "player 6",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_7",
				Description: "player 7",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_8",
				Description: "player 8",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_9",
				Description: "player 9",
				Required:    true,
			},
		},
	}

	ShuffleCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("ShuffleCommandHandler: shuffle started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))

		if err := (*app.KungdotaService).ShufflePlayers(context.Background(), getNames(i)); err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleCommandHandler ShufflePlayers failed %s", err))
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: createResponseData((*app.KungdotaService).GetProperties()),
		}); err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleCommandHandler: InteractionRespond failed %s", err))
		}
	}
)

func getNames(i *discordgo.InteractionCreate) []string {
	oAll := i.ApplicationCommandData().Options
	r := make([]string, 0)

	for _, o := range oAll {
		r = append(r, o.Value.(string))
	}

	return r
}

func createResponseData(p service.Properties) *discordgo.InteractionResponseData {
	btns := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "Move",
			Style:    1,
			Disabled: false,
			CustomID: `move_btn`,
		},
		discordgo.Button{
			Label:    "Update",
			Style:    1,
			Disabled: false,
			CustomID: `update_btn`,
		},
		discordgo.Button{
			Label:    "Cancel",
			Style:    4,
			Disabled: false,
			CustomID: `cancel_btn`,
		},
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Type:  "rich",
			Title: `Team shuffle`,
			Description: fmt.Sprintf(
				"%d\n%s\n%s\n%s\n%s",
				p.ShuffledTeams.EloDiff,
				p.ShuffledTeams.TeamOne.Names(),
				p.ShuffledTeams.TeamTwo.Names(),
				p.ShuffledTeams.FirstPicker.Username,
				p.ShuffledTeams.SecondPicker.Username,
			),
			Color: 0xff00ae,
		},
	}

	return &discordgo.InteractionResponseData{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: btns,
			},
		},
		Embeds: embeds,
	}
}
