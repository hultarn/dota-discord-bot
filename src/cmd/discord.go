package cmd

import "dota-discord-bot/src/cmd/discord"

func init() {
	rootCmd.AddCommand(discord.RootCmd)
}
