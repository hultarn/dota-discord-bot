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

func move() {
	fmt.Println()
}
