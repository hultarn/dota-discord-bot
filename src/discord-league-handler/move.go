package discordleaguehandler

import (
	"context"
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

		teams, err := app.Repository.GetShuffledTeams(context.Background())
		if err != nil {
			app.Logger.Error(fmt.Sprintf("MoveCommandHandler GetShuffledTeam failed %s", err))
			return
		}

		if len(append(teams.TeamOne.Players, teams.TeamTwo.Players...)) < 10 {
			app.Logger.Info("teams not shuffled")
		}
		channelTeamOne, channelTeamTwo := app.DiscordService.GetChannels()
		for _, u := range teams.TeamOne.Players {
			if err := s.GuildMemberMove(i.GuildID, u.DiscordID, &channelTeamOne); err != nil {
				app.Logger.Error(fmt.Sprintf("%v", err))
				continue
			}
		}
		for _, u := range teams.TeamTwo.Players {
			if err := s.GuildMemberMove(i.GuildID, u.DiscordID, &channelTeamTwo); err != nil {
				app.Logger.Error(fmt.Sprintf("%v", err))
				continue
			}
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{}); err != nil {
			return
		}
	}
)
