package application

import (
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		&UpdateCmd,
		&ShuffleCmd,
		&ShuffleSmartCmd,
		&AddGameCmd,
		&MoveCmd,
		&NewPollManualCmd,
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, app application){
		UpdateCmd.Name:        UpdateCommandHandler,
		ShuffleCmd.Name:       ShuffleCommandHandler,
		ShuffleSmartCmd.Name:  ShuffleSmartCommandHandler,
		AddGameCmd.Name:       AddGameCommandHandler,
		MoveCmd.Name:          MoveCommandHandler,
		NewPollManualCmd.Name: NewPollManualHandler,
	}
)
