package discordleaguehandler

import (
	"dota-discord-bot/src/internal/opendota"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	AddGameCmd = discordgo.ApplicationCommand{
		Name:        "addgame",
		Description: "add specified game",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "game-id",
				Description: "game id",
				Required:    true,
			},
		},
	}

	AddGameCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("AddGameCommandHandler: addgame started by user: %s", i.Member.User.Username))

		id := (i.ApplicationCommandData().Options[0].Value).(string)

		var game opendota.OpenDotaGameObject
		game, err := app.OpendotaService.GetMatch(id)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler GetGame failed %s", err))
			return
		}

		if len(game.Objectives) == 0 {
			app.Logger.Info(fmt.Sprintf("match was not parsed, starting parse on match %s", id))
			go func() {
				if err := app.OpendotaService.RequestMatch(id); err != nil {
					app.Logger.Error(fmt.Sprintf("AddGameCommandHandler RequestMatch failed %s", err))
					return
				}
				game, err = app.OpendotaService.GetMatch(id)
				if err != nil {
					app.Logger.Error(fmt.Sprintf("AddGameCommandHandler GetMatch failed %s", err))
					return
				}
			}()
		}

		if err := app.KungdotaService.PostMatch(game); err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler PostMatch failed %s", err))
			return
		}

		if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{}}); err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler InteractionResponseData failed %s", err))
			return
		}
	}
)
