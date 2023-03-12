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
		//TODO: Fix
		tmp := i.ApplicationCommandData().Options
		id := (tmp[0].Value).(string)

		g, err := (*app.OpendotaService).GetMatch(id)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler GetGame failed %s", err))
		}

		if err := (*app.KungdotaService).PostMatch(g); err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler PostMatch failed %s", err))
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{})
	}
)
