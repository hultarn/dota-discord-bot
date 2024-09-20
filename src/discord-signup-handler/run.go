package discordsignuphandler

import (
	db "dota-discord-bot/src/internal/database"
	"dota-discord-bot/src/internal/discord"
	kungdotasvc "dota-discord-bot/src/internal/kungdota/service"
	opendotasvc "dota-discord-bot/src/internal/opendota/service"
	steamdotasvc "dota-discord-bot/src/internal/steamdota/service"
	"os"
	"os/signal"

	"github.com/robfig/cron"
	"go.uber.org/zap"
)

type App struct {
	Logger           *zap.Logger
	DiscordService   DiscordService
	KungdotaService  kungdotasvc.KungdotaService
	SteamdotaService steamdotasvc.SteamdotaService
	OpendotaService  opendotasvc.OpendotaService
	Repository       db.Repository
}

func NewApplication(
	logger *zap.Logger,
	discordService DiscordService,
	kungdotaService kungdotasvc.KungdotaService,
	steamdotaService steamdotasvc.SteamdotaService,
	opendotaService opendotasvc.OpendotaService,
	repository db.Repository,
) discord.Application {
	return &App{
		Logger:           logger,
		DiscordService:   discordService,
		KungdotaService:  kungdotaService,
		SteamdotaService: steamdotaService,
		OpendotaService:  opendotaService,
		Repository:       repository,
	}
}

func (rx *App) Run() {
	if err := rx.DiscordService.SignUpStart(rx); err != nil {
		rx.Logger.Error("Run failed")
		return
	}

	// defer rx.DiscordService.GetProperties().S.Close()

	c := cron.New()
	if err := c.AddFunc("0 0 5 * * 3", func() {
		if err := rx.DiscordService.PostNewSignUpMessage(rx); err != nil {
			rx.Logger.Error("Run failed")
			return
		}

		rx.Logger.Info("Cron PostNewSignUpMessage")
	}); err != nil {
		rx.Logger.Error("Run failed")
		return
	}

	c.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
