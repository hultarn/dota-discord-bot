package application

import (
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		&UpdateCmd,
		&ShuffleCmd,
		&AddGameCmd,
		&MoveCmd,
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, app application){
		UpdateCmd.Name:  UpdateCommandHandler,
		ShuffleCmd.Name: ShuffleCommandHandler,
		AddGameCmd.Name: AddGameCommandHandler,
		MoveCmd.Name:    MoveCommandHandler,
	}
)
