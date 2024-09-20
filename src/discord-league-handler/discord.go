package discordleaguehandler

import (
	"context"
	"dota-discord-bot/src"
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var (
	GuildID = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
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
	PostNewSignUpMessage(app *application) error
	GetChannels() (string, string)
}

type discordService struct {
	config         config
	logger         *zap.Logger
	Session        *discordgo.Session
	SignUpChannel  string
	TeamOneChannel string
	TeamTwoChannel string
}

type config struct {
	token     string
	SignUp    string
	tokenType string
	teamOne   string
	teamTwo   string
}

func NewConfig(token string, tokenType string, signUp string, teamOne string, teamTwo string) config {
	return config{
		token:     token,
		tokenType: tokenType,
		SignUp:    signUp,
		teamOne:   teamOne,
		teamTwo:   teamTwo,
	}
}

func NewDiscordService(logger *zap.Logger, config config) DiscordService {
	r := &discordService{
		config: config,
		logger: logger,
	}

	r.SignUpChannel = config.SignUp
	r.TeamOneChannel = config.teamOne
	r.TeamTwoChannel = config.teamTwo

	return r
}

func (rx *discordService) SignUpStart(app *application) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", rx.config.tokenType, rx.config.token))
	if err != nil {
		panic(err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		rx.logger.Info(fmt.Sprintf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	return nil
}

func (rx *discordService) PostNewSignUpMessage(app *application) error {
	m, err := rx.Session.ChannelMessageSendComplex(rx.SignUpChannel, src.CreateResponseDataSignup())
	if err != nil {
		return err
	}

	if err := app.Repository.InsertCurrentWeekAndYear(context.Background(), m.ID); err != nil {
		return err
	}

	return nil
}

func (rx *discordService) GetChannels() (string, string) {
	return rx.config.teamOne, rx.config.teamTwo
}

func (rx *discordService) Start(app *application) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", rx.config.tokenType, rx.config.token))
	if err != nil {
		panic(err)
	}

	rx.Session = s

	return nil
}

func (rx *discordService) AddHandlers(app *application) error {
	rx.logger.Info("adding Ready event")
	rx.Session.AddHandler(rx.signInHandler())

	if err := rx.Session.Open(); err != nil {
		return err
	}

	rx.logger.Info("adding InteractionCreate event")
	rx.Session.AddHandler(rx.buttonHandler(app))

	return nil
}

func (rx *discordService) AddCommands(app *application) error {
	rx.logger.Info("adding commands...")

	for _, v := range Commands {
		cmd, err := rx.Session.ApplicationCommandCreate(rx.Session.State.User.ID, *GuildID, v)
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

	registeredCommands, err := rx.Session.ApplicationCommands(rx.Session.State.User.ID, *GuildID)
	if err != nil {
		rx.logger.Info(fmt.Sprintf("failed fetching commands: %v", err))
		return err
	}

	for _, v := range registeredCommands {
		if err := rx.Session.ApplicationCommandDelete(rx.Session.State.User.ID, *GuildID, v.ID); err != nil {
			rx.logger.Error(fmt.Sprintf("failed to remove '%v' command: %v", v.Name, err))
			return err
		}
	}

	return nil
}
