package application

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	UpdateCmd = discordgo.ApplicationCommand{
		Name:        "update",
		Description: "fetches new games and updates",
	}

	UpdateCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("UpdateCommandHandler: update started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))

		id, err := (*app.SteamdotaService).GetLatestGameID()
		if err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler GetLatestGameID failed %s", err))
			return
		}

		m, err := (*app.OpendotaService).GetMatch(id)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler GetGame failed %s", err))
			return
		}

		if len(m.Objectives) == 0 {
			app.Logger.Info(fmt.Sprintf("match was not parsed, starting parse on match %s", id))
			go func() {
				if err := (*app.OpendotaService).RequestMatch(id); err != nil {
					app.Logger.Error(fmt.Sprintf("UpdateCommandHandler RequestMatch failed %s", err))
					return
				}
				mParsed, err := (*app.OpendotaService).GetMatch(id)
				if err != nil {
					app.Logger.Error(fmt.Sprintf("UpdateCommandHandler GetMatch failed %s", err))
					return
				}

				if err := (*app.KungdotaService).PostMatch(mParsed); err != nil {
					app.Logger.Error(fmt.Sprintf("UpdateCommandHandler PostMatch failed %s", err))
					return
				}
			}()
		}

		if err := (*app.KungdotaService).PostMatch(m); err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler PostMatch failed %s", err))
			return
		}

		if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{}); err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler PostMatch failed %s", err))
			return
		}
	}
)
