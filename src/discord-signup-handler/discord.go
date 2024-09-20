package discordsignuphandler

import (
	"context"
	"dota-discord-bot/src"
	"dota-discord-bot/src/internal/database"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var buttons = map[string]int{
	"game_1_btn": 0,
	"game_2_btn": 1,
	"game_3_btn": 2,
}

type DiscordService interface {
	SignUpStart(app *App) error
	PostNewSignUpMessage(app *App) error
}

type discordService struct {
	config        config
	logger        *zap.Logger
	Session       *discordgo.Session
	SignUpChannel string
}

type config struct {
	token         string
	SignUpChannel string
	tokenType     string
}

func NewConfig(token string, tokenType string, signUpChannel string) config {
	return config{
		token:         token,
		tokenType:     tokenType,
		SignUpChannel: signUpChannel,
	}
}

func NewDiscordService(logger *zap.Logger, config config) DiscordService {
	r := &discordService{
		config: config,
		logger: logger,
	}

	r.SignUpChannel = config.SignUpChannel

	return r
}

func (rx *discordService) PostNewSignUpMessage(app *App) error {
	m, err := rx.Session.ChannelMessageSendComplex(rx.SignUpChannel, src.CreateResponseDataSignup())
	if err != nil {
		return err
	}

	if err := app.Repository.InsertCurrentWeekAndYear(context.Background(), m.ID); err != nil {
		return err
	}

	return nil
}

func format(list []database.Sign) (s []string) {
	s = make([]string, 0)
	for i := range list {
		s = append(s, list[i].DiscordID)
	}

	return
}

func (rx *discordService) SignUpStart(app *App) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", rx.config.tokenType, rx.config.token))
	if err != nil {
		panic(err)
	}

	s.AddHandler(rx.signInHandler())
	s.AddHandler(rx.signUpHandler(app))

	if err := s.Open(); err != nil {
		return err
	}

	rx.Session = s

	return nil
}
