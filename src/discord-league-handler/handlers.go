package discordleaguehandler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (rx *discordService) signInHandler() func(s *discordgo.Session, i *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		rx.logger.Info(fmt.Sprintf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	}
}

func (rx *discordService) buttonHandler(app *application) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch v := i.Type; v {
		case discordgo.InteractionApplicationCommand:
			if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i, *app)
			}
		case discordgo.InteractionMessageComponent:
			switch v := i.MessageComponentData().CustomID; v {
			case moveBtn:
				MoveCommandHandler(s, i, *app)
			case updateBtn:
				UpdateCommandHandler(s, i, *app)
			case cancelBtn:
				if err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID); err != nil {
					return
				}
			}
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{}); err != nil {
				return
			}
		}
	}
}
