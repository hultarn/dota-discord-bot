package discordleaguehandler

import (
	db "dota-discord-bot/src/internal/database"
	"dota-discord-bot/src/internal/discord"
	kungdotasvc "dota-discord-bot/src/internal/kungdota/service"
	opendotasvc "dota-discord-bot/src/internal/opendota/service"
	steamdotasvc "dota-discord-bot/src/internal/steamdota/service"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

type application struct {
	Logger           *zap.Logger
	DiscordService   DiscordService
	KungdotaService  kungdotasvc.KungdotaService
	SteamdotaService steamdotasvc.SteamdotaService
	OpendotaService  opendotasvc.OpendotaService
	Repository       db.Repository
	Misc             *Misc
}

type Misc struct {
	Admin []string
}

func NewApplication(
	logger *zap.Logger,
	discordService DiscordService,
	kungdotaService kungdotasvc.KungdotaService,
	steamdotaService steamdotasvc.SteamdotaService,
	opendotaService opendotasvc.OpendotaService,
	repository db.Repository,
	misc *Misc,
) discord.Application {
	return &application{
		Logger:           logger,
		DiscordService:   discordService,
		KungdotaService:  kungdotaService,
		SteamdotaService: steamdotaService,
		OpendotaService:  opendotaService,
		Repository:       repository,
		Misc:             misc,
	}
}

func (rx *application) Run() {
	if err := rx.DiscordService.Start(rx); err != nil {
		rx.Logger.Error("Run failed")
		return
	}

	if err := rx.DiscordService.AddHandlers(rx); err != nil {
		rx.Logger.Error("Run failed")
		return
	}
	if err := rx.DiscordService.AddCommands(rx); err != nil {
		rx.Logger.Error("Run failed")
		return
	}
	// rx.DiscordService.RemoveCommands(rx)
	// defer rx.DiscordService.GetProperties().S.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// if err := rx.DiscordService.RemoveCommands(rx); err != nil {
	// 	rx.Logger.Error("Run failed")
	// 	return
	// }
}
