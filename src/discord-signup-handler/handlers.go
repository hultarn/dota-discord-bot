package discordsignuphandler

import (
	"context"
	"dota-discord-bot/src"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func (rx *discordService) signInHandler() func(s *discordgo.Session, i *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		rx.logger.Info(fmt.Sprintf("logged in as: %v", s.State.User.Username))
	}
}

func (rx *discordService) signUpHandler(app *App) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ctx := context.Background()

		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		id, err := app.Repository.GetByCurrentWeekAndYear(ctx)
		if err != nil {
			return
		}

		if i.Message.ID != id {
			return
		}

		// To avoid red warning in discord.
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		}); err != nil {
			rx.logger.Info(fmt.Sprintf("InteractionRespond %v", err))
			return
		}

		_, err = app.KungdotaService.GetPlayersByDiscordIDs(ctx, []string{i.Interaction.Member.User.ID})
		if err != nil {
			rx.logger.Info(fmt.Sprintf("Player %s doesn't exist", i.Interaction.Member.User.Username))

			// TODO Fix automated signup:
			ch, err := s.UserChannelCreate(i.Interaction.Member.User.ID)
			if err != nil {
				rx.logger.Error(fmt.Sprintln(err))
			}

			_, err = s.ChannelMessageSend(ch.ID, "Seems like you don't exist.. you should probably talk to someone.")
			if err != nil {
				rx.logger.Error(fmt.Sprintln(err))
			}

			return
		}

		if err := app.Repository.InsertPlayer(ctx, i.Interaction.Member.User.ID, buttons[i.MessageComponentData().CustomID]); err != nil {
			rx.logger.Error("failed InsertPlayer")
			return
		}

		entry, err := app.Repository.GetCurrent(ctx)
		if err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}

		embeds := src.GetEmbeds()
		for i, game := range []uuid.UUID{entry.Game_1, entry.Game_2, entry.Game_3} {
			players, err := app.Repository.GetSigedPlayersByGame(ctx, game.String())
			if err != nil {
				rx.logger.Error(fmt.Sprintln(err))
			}

			d, err := app.KungdotaService.GetPlayersByDiscordIDs(ctx, format(players))
			if err != nil {
				rx.logger.Error(fmt.Sprintln(err))
			}

			embeds[i].Description = strings.Join(d, ", ") + " Tot:" + fmt.Sprint(len(d))
		}

		if _, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embeds:  &embeds,
			ID:      id,
			Channel: rx.SignUpChannel,
		}); err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}
	}
}
