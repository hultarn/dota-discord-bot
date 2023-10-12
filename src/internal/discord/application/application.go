package application

import (
	dynamodbsvc "dota-discord-bot/src/internal/dynamodb/service"
	kungdotasvc "dota-discord-bot/src/internal/kungdota/service"
	opendotasvc "dota-discord-bot/src/internal/opendota/service"
	steamdotasvc "dota-discord-bot/src/internal/steamdota/service"
	"flag"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

var (
	GuildID = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
)

type Application interface {
	Run()
	RunSignUp()
}

type Misc struct {
	SuperDuperAdmin []string
}
type application struct {
	Logger           *zap.Logger
	DiscordService   *DiscordService
	KungdotaService  *kungdotasvc.KungdotaService
	SteamdotaService *steamdotasvc.SteamdotaService
	OpendotaService  *opendotasvc.OpendotaService
	DynamodbService  *dynamodbsvc.DynamodbService
	Misc             *Misc
}

func NewApplication(
	logger *zap.Logger,
	discordService *DiscordService,
	kungdotaService *kungdotasvc.KungdotaService,
	steamdotaService *steamdotasvc.SteamdotaService,
	opendotaService *opendotasvc.OpendotaService,
	dynamodbService *dynamodbsvc.DynamodbService,
	misc *Misc,
) Application {
	return &application{
		Logger:           logger,
		DiscordService:   discordService,
		KungdotaService:  kungdotaService,
		SteamdotaService: steamdotaService,
		OpendotaService:  opendotaService,
		DynamodbService:  dynamodbService,
		Misc:             misc,
	}
}

func (rx *application) RunSignUp() {
	(*rx.DiscordService).SignUpStart(rx)

	defer (*rx.DiscordService).GetProperties().S.Close()

	// TODO: Better way?
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func (rx *application) Run() {
	(*rx.DiscordService).Start(rx)
	(*rx.DiscordService).AddHandlers(rx)
	(*rx.DiscordService).AddCommands(rx)

	defer (*rx.DiscordService).GetProperties().S.Close()

	// TODO: Better way?
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	(*rx.DiscordService).RemoveCommands(rx)
}
