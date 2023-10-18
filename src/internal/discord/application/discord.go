package application

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

const (
	moveBtn   = "move_btn"
	updateBtn = "update_btn"
	cancelBtn = "cancel_btn"

	gameOneBtn   = "game_1_btn"
	gameTwoBtn   = "game_2_btn"
	gameThreeBtn = "game_3_btn"
	gameClearBtn = "game_clear_btn"
)

type DiscordService interface {
	Start(app *application) error
	AddHandlers(app *application) error
	AddCommands(app *application) error
	RemoveCommands(app *application) error
	GetProperties() properties

	SignUpStart(app *application) error
	PostNewSignUpMessage(app *application) error
}

type discordService struct {
	config config
	logger *zap.Logger
	props  properties
}

type config struct {
	token     string
	SignUp    string
	tokenType string
	teamOne   string
	teamTwo   string
}

type properties struct {
	S       *discordgo.Session
	SignUp  string
	TeamOne string
	TeamTwo string
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

	r.props.SignUp = config.SignUp
	r.props.TeamOne = config.teamOne
	r.props.TeamTwo = config.teamTwo

	return r
}

func (rx *discordService) PostNewSignUpMessage(app *application) error {
	m, err := rx.props.S.ChannelMessageSendComplex((*app.DiscordService).GetProperties().SignUp, createResponseDataSignup())
	if err != nil {
		return err
	}

	if err := (*app.DynamodbService).InsertCurrentWeekAndYear(context.Background(), m.ID); err != nil {
		return err
	}

	return nil
}

func (rx *discordService) SignUpStart(app *application) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", rx.config.tokenType, rx.config.token))
	if err != nil {
		panic(err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		rx.logger.Info(fmt.Sprintf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		id, err := (*app.DynamodbService).GetByCurrentWeekAndYear(context.Background())
		if err != nil {
			return
		}

		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		if i.Message.ID != id {
			return
		}

		id, err = (*app.DynamodbService).GetByCurrentWeekAndYear(context.Background())
		if err != nil {

			return
		}

		if id == "" {
			m, err := s.ChannelMessageSendComplex((*app.DiscordService).GetProperties().SignUp, createResponseDataSignup())
			if err != nil {
				return
			}

			if err := (*app.DynamodbService).InsertCurrentWeekAndYear(context.Background(), m.ID); err != nil {
				return
			}
		}

		_, err = (*app.KungdotaService).GetPlayersByDiscordIDs(context.Background(), []string{i.Interaction.Member.User.ID})
		if err != nil {
			rx.logger.Info(fmt.Sprintf("Player %s doesn't exist", i.Interaction.Member.User.Username))

			// TOOD Fix automated signup:
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

		switch v := i.MessageComponentData().CustomID; v {
		case gameOneBtn:
			//newList, err = (*app.KungdotaService).SignUp(context.Background(), i.Interaction.Member.User.Username, 0)
			if err := (*app.DynamodbService).InsertPlayer(context.Background(), i.Interaction.Member.User.ID, 0); err != nil {
				rx.logger.Error("failed InsertPlayer")
			}
		case gameTwoBtn:
			// newList, err = (*app.KungdotaService).SignUp(context.Background(), i.Interaction.Member.User.Username, 1)
			if err := (*app.DynamodbService).InsertPlayer(context.Background(), i.Interaction.Member.User.ID, 1); err != nil {
				rx.logger.Error("failed InsertPlayer")
			}
		case gameThreeBtn:
			// newList, err = (*app.KungdotaService).SignUp(context.Background(), i.Interaction.Member.User.Username, 2)
			if err := (*app.DynamodbService).InsertPlayer(context.Background(), i.Interaction.Member.User.ID, 2); err != nil {
				rx.logger.Error("failed InsertPlayer")
			}
		case gameClearBtn:
			if slices.Contains(app.Misc.SuperDuperAdmin, i.Interaction.Member.User.ID) {
				if err := (*app.DynamodbService).ClearPlayers(context.Background()); err != nil {
					rx.logger.Error("failed InsertPlayer")
				}
			} else {
				ch, err := s.UserChannelCreate(i.Interaction.Member.User.ID)
				if err != nil {
					rx.logger.Error(fmt.Sprintln(err))
				}

				_, err = s.ChannelMessageSend(ch.ID, "Ät skit.")
				if err != nil {
					rx.logger.Error(fmt.Sprintln(err))
				}
			}
		}

		// if err != nil {
		// 	rx.logger.Error(fmt.Sprintf("failed to signup %v", err))
		// 	return
		// }

		e, err := (*app.DynamodbService).GetCurrent(context.Background())
		if err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}

		d1, err := (*app.KungdotaService).GetPlayersByDiscordIDs(context.Background(), e.Game_1)
		if err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}

		d2, err := (*app.KungdotaService).GetPlayersByDiscordIDs(context.Background(), e.Game_2)
		if err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}

		d3, err := (*app.KungdotaService).GetPlayersByDiscordIDs(context.Background(), e.Game_3)
		if err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}

		embeds := getEmbeds()
		embeds[0].Description = strings.Join(d1, ", ") + " Tot:" + fmt.Sprint(len(d1))
		embeds[1].Description = strings.Join(d2, ", ") + " Tot:" + fmt.Sprint(len(d2))
		embeds[2].Description = strings.Join(d3, ", ") + " Tot:" + fmt.Sprint(len(d3))

		if _, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embeds:  embeds,
			ID:      id,
			Channel: (*app.DiscordService).GetProperties().SignUp,
		}); err != nil {
			rx.logger.Error(fmt.Sprintln(err))
		}
	})

	if err := s.Open(); err != nil {
		return err
	}

	rx.props.S = s

	return nil
}

func createResponseDataSignup() *discordgo.MessageSend {
	btns := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "Game_1",
			Style:    1,
			Disabled: false,
			CustomID: gameOneBtn,
		},
		discordgo.Button{
			Label:    "Game_2",
			Style:    1,
			Disabled: false,
			CustomID: gameTwoBtn,
		},
		discordgo.Button{
			Label:    "Game_3",
			Style:    1,
			Disabled: false,
			CustomID: gameThreeBtn,
		},
		discordgo.Button{
			Label:    "Clear",
			Style:    4,
			Disabled: false,
			CustomID: gameClearBtn,
		},
	}

	embeds := getEmbeds()

	return &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: btns,
			},
		},
		Embeds: embeds,
	}
}

func getEmbeds() []*discordgo.MessageEmbed {
	// TODO add times from .env
	return []*discordgo.MessageEmbed{
		{
			Type:  "rich",
			Title: `Game_1: 19:30`,
			Color: 0xff00ae,
		},
		{
			Type:  "rich",
			Title: `Game_2: 20:45`,
			Color: 0xff00ae,
		},
		{
			Type:  "rich",
			Title: `Game_3: 22:00`,
			Color: 0xff00ae,
		},
	}
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
				if err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID); err != nil {
					return
				}
			}
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{}); err != nil {
				return
			}
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
