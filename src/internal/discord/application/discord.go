package application

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

const (
	moveBtn   = "move_btn"
	updateBtn = "update_btn"
	cancelBtn = "cancel_btn"

	gameOneBtn    = "game_1_btn"
	gameTwoBtn    = "game_2_btn"
	gameThreeBtn  = "game_3_btn"
	gameClearBtn  = "game_clear_btn"
	gameUpdateBtn = "game_update_btn"
)

type DiscordService interface {
	Start(app *application) error
	AddHandlers(app *application) error
	AddCommands(app *application) error
	RemoveCommands(app *application) error
	GetProperties() properties

	SignUpStart(app *application) error
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

func (rx *discordService) SignUpStart(app *application) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", rx.config.tokenType, rx.config.token))
	if err != nil {
		panic(err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		rx.logger.Info(fmt.Sprintf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))

		id, err := (*app.DynamodbService).GetByCurrentWeekAndYear(context.Background())
		if err != nil {
			return
		}
		if id == "" {
			m, err := s.ChannelMessageSendComplex("801048845055426560", createResponseData2())
			if err != nil {
				return
			}

			if err := (*app.DynamodbService).InsertCurrentWeekAndYear(context.Background(), m.ID); err != nil {
				return
			}
		}
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

		embeds := []*discordgo.MessageEmbed{
			{
				Type:        "rich",
				Title:       `Time                19:30    20:45    22:00`,
				Description: fmt.Sprintf("%d\n%d\n%d", 1, 2, 3),
				Color:       0xff00ae,
			},
		}

		var newList map[string][]string

		switch v := i.MessageComponentData().CustomID; v {
		case gameOneBtn:
			newList, _ = (*app.KungdotaService).SignUp(context.Background(), i.Interaction.Member.User.Username, 0)
			// (*app.DynamodbService).
		case gameTwoBtn:
			newList, _ = (*app.KungdotaService).SignUp(context.Background(), i.Interaction.Member.User.Username, 1)
		case gameThreeBtn:
			newList, _ = (*app.KungdotaService).SignUp(context.Background(), i.Interaction.Member.User.Username, 2)
		case gameClearBtn:

		case gameUpdateBtn:
			newList, _ = (*app.KungdotaService).Update(context.Background(), i.Interaction.Member.User.Username)
			// (*app.DynamodbService).
		}

		tmp := ""
		for s, k := range newList {
			tmp += s + ": "
			for _, k := range k {
				tmp += k + ", "
			}
			tmp += "\n"
			tmp += fmt.Sprintf("%d", len(k))
			tmp += "\n"
		}

		embeds[0].Description = tmp

		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embeds:  embeds,
			ID:      id,
			Channel: "801048845055426560",
		})

	})

	if err := s.Open(); err != nil {
		return err
	}

	rx.props.S = s

	return nil
}

func createResponseData2() *discordgo.MessageSend {
	btns := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "Game_1",
			Style:    1,
			Disabled: false,
			CustomID: `game_1_btn`,
		},
		discordgo.Button{
			Label:    "Game_2",
			Style:    1,
			Disabled: false,
			CustomID: `game_2_btn`,
		},
		discordgo.Button{
			Label:    "Game_3",
			Style:    1,
			Disabled: false,
			CustomID: `game_3_btn`,
		},
		discordgo.Button{
			Label:    "Clear",
			Style:    4,
			Disabled: false,
			CustomID: `game_clear_btn`,
		},
		discordgo.Button{
			Label:    "Update",
			Style:    4,
			Disabled: false,
			CustomID: `game_update_btn`,
		},
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Type:        "rich",
			Title:       `Time                19:30    20:45    22:00`,
			Description: fmt.Sprintf("%d\n%d\n%d", 1, 2, 3),
			Color:       0xff00ae,
		},
	}

	return &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: btns,
			},
		},
		Embeds: embeds,
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
