package discord

import (
	"dota-discord-bot/src/internal/discord/application"
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
	Use: "start-discord-bot",
	Run: startDiscordBot,
}

var app application.Application

func init() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	dsvc := application.NewDiscordService(logger,
		application.NewConfig(config.Token, "Bot", config.TeamOne, config.TeamTwo),
	)

	httpC := http.Client{}

	ksvc := kungdotasvc.NewKungdotaService(logger,
		kungdotarepo.NewKungdotaRepository(logger, httpC,
			kungdotarepo.NewConfig(config.KungDotaID)),
	)

	ssvc := steamdotasvc.NewSteamdotaService(logger,
		steamdotarepo.NewSteamdotaRepository(logger, httpC,
			steamdotarepo.NewConfig(config.DotaID, config.SteamKey)),
	)

	osvc := opendotasvc.NewOpendotaService(logger,
		opendotarepo.NewOpendotaRepository(logger, httpC,
			opendotarepo.NewConfig(config.SteamKey)),
	)

	app = application.NewApplication(logger, &dsvc, &ksvc, &ssvc, &osvc)
}

func startDiscordBot(cmd *cobra.Command, _ []string) {
	app.Run()
}

func init() {
	RootCmd.AddCommand(startDiscordBotCmd)
}
