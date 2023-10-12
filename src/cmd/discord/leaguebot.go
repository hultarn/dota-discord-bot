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
	Use: "start-league-bot",
	Run: startLeagueBot,
}

var leagueBot application.Application

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

	leagueBot = application.NewApplication(logger, &dsvc, &ksvc, &ssvc, &osvc, nil, nil)
}

func startLeagueBot(cmd *cobra.Command, _ []string) {
	leagueBot.Run()
}

func init() {
	RootCmd.AddCommand(startDiscordBotCmd)
}
