package application

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	MoveCmd = discordgo.ApplicationCommand{
		Name:        "move",
		Description: "moves players after shuffle",
	}

	MoveCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("MoveCommandHandler: move started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))

		props := (*app.KungdotaService).GetProperties().ShuffledTeams

		if len(append(props.TeamOne.Players, props.TeamTwo.Players...)) < 10 {
			app.Logger.Info("teams not shuffled")
		}

		tOne := (*app.DiscordService).GetProperties().TeamOne
		tTwo := (*app.DiscordService).GetProperties().TeamTwo
		for _, u := range props.TeamOne.Players {
			if err := s.GuildMemberMove(i.GuildID, u.DiscordID, &tOne); err != nil {
				app.Logger.Error(fmt.Sprintf("%v", err))
				continue
			}
		}
		for _, u := range props.TeamTwo.Players {
			if err := s.GuildMemberMove(i.GuildID, u.DiscordID, &tTwo); err != nil {
				app.Logger.Error(fmt.Sprintf("%v", err))
				continue
			}
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{}); err != nil {
			return
		}
	}
)
