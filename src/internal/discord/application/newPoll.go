package application

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	NewPollManualCmd = discordgo.ApplicationCommand{
		Name:        "newpoll",
		Description: "shuffle teams from signup list",
	}

	NewPollManualHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("NewPollManualHandler: new poll started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))
		if err := (*app.DiscordService).PostNewSignUpMessage(&app); err != nil {
			app.Logger.Error(fmt.Sprintf("NewPollManualHandler PostNewSignUpMessage failed %s", err))
			return
		}
	}
)
