package discord

import (
	discordsignuphandler "dota-discord-bot/src/discord-signup-handler"
	db "dota-discord-bot/src/internal/database/mysql"
	"dota-discord-bot/src/internal/discord"
	kungdotarepo "dota-discord-bot/src/internal/kungdota/repository"
	kungdotasvc "dota-discord-bot/src/internal/kungdota/service"
	opendotarepo "dota-discord-bot/src/internal/opendota/repository"
	opendotasvc "dota-discord-bot/src/internal/opendota/service"
	steamdotarepo "dota-discord-bot/src/internal/steamdota/repository"
	steamdotasvc "dota-discord-bot/src/internal/steamdota/service"

	"net/http"
	"net/http/cookiejar"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var startSignupCmd = &cobra.Command{
	Use: "start-signup-bot",
	Run: startSignupBot,
}

var signupBot discord.Application
var envPathSignup string

func init() {
	startSignupCmd.Flags().StringVarP(&envPathSignup, "env", "e", "", "Path to the .env file")
	RootCmd.AddCommand(startSignupCmd)
}

func initSignupBot() {
	config, err := NewConfig(envPathSignup)
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	dsvc := discordsignuphandler.NewDiscordService(logger,
		discordsignuphandler.NewConfig(config.Token, "Bot", config.SignUp),
	)

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpC := http.Client{Jar: jar}

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

	signupBot = discordsignuphandler.NewApplication(logger, dsvc, ksvc, ssvc, osvc, dbr)
}

func startSignupBot(cmd *cobra.Command, _ []string) {
	initSignupBot()
	signupBot.Run()
}
