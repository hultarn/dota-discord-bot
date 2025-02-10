package discord

import (
	discordleaguehandler "dota-discord-bot/src/discord-league-handler"
	db "dota-discord-bot/src/internal/database/mysql"
	"dota-discord-bot/src/internal/discord"
	kungdotarepo "dota-discord-bot/src/internal/kungdota/repository"
	kungdotasvc "dota-discord-bot/src/internal/kungdota/service"
	opendotarepo "dota-discord-bot/src/internal/opendota/repository"
	opendotasvc "dota-discord-bot/src/internal/opendota/service"
	steamdotarepo "dota-discord-bot/src/internal/steamdota/repository"
	steamdotasvc "dota-discord-bot/src/internal/steamdota/service"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var startDiscordBotCmd = &cobra.Command{
	Use: "start-league-bot",
	Run: startLeagueBot,
}

var leagueBot discord.Application
var envPathLeague string

func init() {
	startDiscordBotCmd.Flags().StringVarP(&envPathLeague, "env", "e", "", "Path to the .env file")
	RootCmd.AddCommand(startDiscordBotCmd)
}

func initLeagueBot() {
	config, err := NewConfig(envPathLeague)
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	dsvc := discordleaguehandler.NewDiscordService(logger,
		discordleaguehandler.NewConfig(config.Token, "Bot", config.SignUp, config.TeamOne, config.TeamTwo),
	)

	httpC := http.Client{}

	ksvc := kungdotasvc.NewKungdotaService(logger,
		kungdotarepo.NewKungdotaRepository(logger, httpC,
			kungdotarepo.NewConfig(config.KungDotaID)),
		config.KungDotaID,
	)

	ssvc := steamdotasvc.NewSteamdotaService(logger,
		steamdotarepo.NewSteamdotaRepository(logger, httpC,
			steamdotarepo.NewConfig(config.SteamKey, config.DotaID)),
	)

	osvc := opendotasvc.NewOpendotaService(logger,
		opendotarepo.NewOpendotaRepository(logger, httpC,
			opendotarepo.NewConfig(config.SteamKey)),
	)

	dbr := db.NewMySqlRepository(logger, db.NewConfig(config.MySqlString))

	leagueBot = discordleaguehandler.NewApplication(logger, dsvc, ksvc, ssvc, osvc, dbr, nil)
}

func startLeagueBot(cmd *cobra.Command, _ []string) {
	initLeagueBot()
	leagueBot.Run()
}
