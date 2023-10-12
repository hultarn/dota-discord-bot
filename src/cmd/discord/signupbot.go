package discord

import (
	"context"
	"dota-discord-bot/src/internal/discord/application"
	dynamodbrepo "dota-discord-bot/src/internal/dynamodb/repository"
	dynamodbsvc "dota-discord-bot/src/internal/dynamodb/service"
	kungdotarepo "dota-discord-bot/src/internal/kungdota/repository"
	kungdotasvc "dota-discord-bot/src/internal/kungdota/service"
	opendotarepo "dota-discord-bot/src/internal/opendota/repository"
	opendotasvc "dota-discord-bot/src/internal/opendota/service"
	steamdotarepo "dota-discord-bot/src/internal/steamdota/repository"
	steamdotasvc "dota-discord-bot/src/internal/steamdota/service"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"net/http"
	"net/http/cookiejar"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var startSignupCmd = &cobra.Command{
	Use: "start-signup-bot",
	Run: startSignupBot,
}

var signupBot application.Application

func init() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	dsvc := application.NewDiscordService(logger,
		application.NewConfig(config.Token, "Bot", config.SignUp, config.TeamOne, config.TeamTwo),
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

	dydbsvc := dynamodbsvc.NewDynamodbService(logger,
		dynamodbrepo.NewDynamodbRepository(logger,
			dynamodbrepo.NewConfig(awsConfig)),
	)

	misc := application.Misc{
		SuperDuperAdmin: config.SuperDuperAdmin,
	}

	signupBot = application.NewApplication(logger, &dsvc, &ksvc, &ssvc, &osvc, &dydbsvc, &misc)
}

func startSignupBot(cmd *cobra.Command, _ []string) {
	signupBot.RunSignUp()
}

func init() {
	RootCmd.AddCommand(startSignupCmd)
}
