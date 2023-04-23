package application

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	NewUserCmd = discordgo.ApplicationCommand{
		Name:        "new_user",
		Description: "new user",
	}

	NewUserCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("NewUserCommandHandler: new_user started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))

		//if todo add not already user check

		props := (*app.KungdotaService).GetProperties().ShuffledTeams

		if len(append(props.TeamOne.Players, props.TeamTwo.Players...)) < 10 {
			app.Logger.Info("teams not shuffled")
		}

		tOne := (*app.DiscordService).GetProperties().TeamOne
		tTwo := (*app.DiscordService).GetProperties().TeamTwo
		for _, u := range props.TeamOne.Players {
			s.GuildMemberMove(i.GuildID, u.DiscordID, &tOne)
		}
		for _, u := range props.TeamTwo.Players {
			s.GuildMemberMove(i.GuildID, u.DiscordID, &tTwo)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{})
	}
)
