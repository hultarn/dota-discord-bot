package discordleaguehandler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var (
	ShuffleSmartCmd = discordgo.ApplicationCommand{
		Name:        "shufflesmart",
		Description: "shuffle teams from signup list",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "i",
				Description: "game index",
				Required:    true,
			},
		},
	}

	ShuffleSmartCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate, app application) {
		app.Logger.Info(fmt.Sprintf("ShuffleSmartCommandHandler: shuffle started by user: %s#%s", i.Member.User.Username, i.Member.User.Discriminator))

		//TODO: Fix
		tmp := i.ApplicationCommandData().Options
		idx := int((tmp[0].Value).(float64))

		//TODO: This is pretty dumb...
		p, err := app.Repository.GetCurrentPlayers(context.Background(), idx)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler ShufflePlayers failed %v", err))
			return
		}

		p2, err := app.KungdotaService.GetPlayersByDiscordIDs(context.Background(), p)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler ShufflePlayers failed %v", err))
			return
		}

		p3, err := app.KungdotaService.GetPlayersByNames(context.Background(), p2)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler GetPlayersByNames failed %v", err))
			return
		}

		if len(p3.Players) != 10 {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler failed invalid amount %d", len(p3.Players)))
			return
		}

		shuffledTeams, err := app.KungdotaService.ShufflePlayers(context.Background(), p3)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler ShufflePlayers failed %s", err))
			return
		}

		if err = app.Repository.InsertShuffledPlayers(context.Background(), shuffledTeams, uuid.New()); err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler InsertShuffledPlayers failed %s", err))
			return
		}

		teams, err := app.Repository.GetShuffledTeams(context.Background())
		if err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler failed invalid amount %d", len(p3.Players)))
			return
		}
		data := createResponseData(teams)

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: data,
		}); err != nil {
			app.Logger.Error(fmt.Sprintf("ShuffleSmartCommandHandler: InteractionRespond failed %s", err))
			return
		}
	}
)
