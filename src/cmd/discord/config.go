package discord

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Token       string
	SignUp      string
	TeamOne     string
	TeamTwo     string
	SteamKey    string
	DotaID      string
	KungDotaID  string
	Admin       []string
	MySqlString string
}

func NewConfig() (Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return Config{}, fmt.Errorf("err executable path")
	}
	fmt.Println(filepath.Dir(ex), ex)
	if err = godotenv.Load(fmt.Sprintf("%s/.env", filepath.Dir(ex))); err != nil {
		return Config{}, fmt.Errorf("err loading .env")
	}

	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	signUp, ok := os.LookupEnv("DISCORD_SIGNUP_ID")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	teamOne, ok := os.LookupEnv("DISCORD_TEAM_ONE_ID")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	teamTwo, ok := os.LookupEnv("DISCORD_TEAM_TWO_ID")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	steamKey, ok := os.LookupEnv("STEAM_API_KEY")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	dotaID, ok := os.LookupEnv("STEAM_DOTA_LEAGUE_ID")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	kungDotaID, ok := os.LookupEnv("KUNGDOTA_LEAGE_ID")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	admin, ok := os.LookupEnv("ADMIN")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	mySqlString, ok := os.LookupEnv("MYSQL_STRING")
	if !ok {
		return Config{}, fmt.Errorf("err loading .env")
	}

	return Config{
		Token:       discordToken,
		SignUp:      signUp,
		TeamOne:     teamOne,
		TeamTwo:     teamTwo,
		SteamKey:    steamKey,
		DotaID:      dotaID,
		KungDotaID:  kungDotaID,
		Admin:       strings.Split(admin, ","),
		MySqlString: mySqlString,
	}, nil
}
