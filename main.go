package main

import (
	"dota-discord-bot/src/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Panic()
		}
	}()

	cmd.Execute()
}
