package application

import (
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
		// TODO: Add permissions?

		app.Logger.Info(fmt.Sprintf("AddGameCommandHandler: addgame started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))

		//TODO: Fix
		tmp := i.ApplicationCommandData().Options
		id := (tmp[0].Value).(string)

		g, err := (*app.OpendotaService).GetMatch(id)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler GetGame failed %s", err))
			return
		}

		if len(g.Objectives) == 0 {
			app.Logger.Info(fmt.Sprintf("match was not parsed, starting parse on match %s", id))
			go func() {
				if err := (*app.OpendotaService).RequestMatch(id); err != nil {
					app.Logger.Error(fmt.Sprintf("AddGameCommandHandler RequestMatch failed %s", err))
					return
				}
				mParsed, err := (*app.OpendotaService).GetMatch(id)
				if err != nil {
					app.Logger.Error(fmt.Sprintf("AddGameCommandHandler GetMatch failed %s", err))
					return
				}

				if err := (*app.KungdotaService).PostMatch(mParsed); err != nil {
					app.Logger.Error(fmt.Sprintf("AddGameCommandHandler PostMatch failed %s", err))
					return
				}

				return
			}()
		} else {
			if err := (*app.KungdotaService).PostMatch(g); err != nil {
				app.Logger.Error(fmt.Sprintf("AddGameCommandHandler PostMatch failed %s", err))
				return
			}
		}

		if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: 4,
			Data: &discordgo.InteractionResponseData{}}); err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler InteractionResponseData failed %s", err))
			return
		}
	}
)
