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

		if err := (*app.KungdotaService).PostMatch(g); err != nil {
			app.Logger.Error(fmt.Sprintf("AddGameCommandHandler PostMatch failed %s", err))
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{})
	}
)
