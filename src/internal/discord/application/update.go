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
		app.Logger.Info(fmt.Sprintf("UpdateCommandHandler: update started by user: %s", i.Member.Nick))

		id, err := (*app.SteamdotaService).GetLatestGameID()
		if err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler GetLatestGameID failed %s", err))
		}

		g, err := (*app.OpendotaService).GetMatch(id)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler GetGame failed %s", err))
		}

		if err := (*app.KungdotaService).PostMatch(g); err != nil {
			app.Logger.Error(fmt.Sprintf("UpdateCommandHandler PostMatch failed %s", err))
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{})
	}
)
