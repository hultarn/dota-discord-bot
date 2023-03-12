package application

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

const (
	moveBtn   = "move_btn"
	updateBtn = "update_btn"
	cancelBtn = "cancel_btn"
)

type DiscordService interface {
	Start(app *application) error
	AddHandlers(app *application) error
	AddCommands(app *application) error
	RemoveCommands(app *application) error
	GetProperties() properties
}

type discordService struct {
	config config
	logger *zap.Logger
	props  properties
}

type config struct {
	token     string
	tokenType string
	teamOne   string
	teamTwo   string
}

type properties struct {
	S       *discordgo.Session
	TeamOne string
	TeamTwo string
}

func NewConfig(token string, tokenType string, teamOne string, teamTwo string) config {
	return config{
		token:     token,
		tokenType: tokenType,
		teamOne:   teamOne,
		teamTwo:   teamTwo,
	}
}

func NewDiscordService(logger *zap.Logger, config config) DiscordService {
	r := &discordService{
		config: config,
		logger: logger,
	}

	r.props.TeamOne = config.teamOne
	r.props.TeamTwo = config.teamTwo

	return r
}

func (rx *discordService) Start(app *application) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", rx.config.tokenType, rx.config.token))
	if err != nil {
		panic(err)
	}

	rx.props.S = s

	return nil
}

func (rx *discordService) AddHandlers(app *application) error {
	rx.logger.Info("adding Ready event")
	rx.props.S.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		rx.logger.Info(fmt.Sprintf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	if err := rx.props.S.Open(); err != nil {
		return err
	}

	rx.logger.Info("adding InteractionCreate event")
	rx.props.S.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{})
		}
	})

	return nil
}

func (rx *discordService) AddCommands(app *application) error {
	rx.logger.Info("adding commands...")

	for _, v := range Commands {
		cmd, err := rx.props.S.ApplicationCommandCreate(rx.props.S.State.User.ID, *GuildID, v)
		if err != nil {
			rx.logger.Error(fmt.Sprintf("failed to create '%v' command: %v", v.Name, err))
			return err
		}
		rx.logger.Info(fmt.Sprintf("added %s", cmd.Name))
	}

	return nil
}

func (rx *discordService) RemoveCommands(app *application) error {
	rx.logger.Info("removing commands...")

	registeredCommands, err := rx.props.S.ApplicationCommands(rx.props.S.State.User.ID, *GuildID)
	if err != nil {
		rx.logger.Info(fmt.Sprintf("failed fetching commands: %v", err))
		return err
	}

	for _, v := range registeredCommands {
		if err := rx.props.S.ApplicationCommandDelete(rx.props.S.State.User.ID, *GuildID, v.ID); err != nil {
			rx.logger.Error(fmt.Sprintf("failed to remove '%v' command: %v", v.Name, err))
			return err
		}
	}

	return nil
}

func (rx *discordService) GetProperties() properties {
	return rx.props
}
