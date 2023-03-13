package repository

import (
	"dota-discord-bot/src/internal/steamdota"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

const endpoint = "https://api.steampowered.com/IDOTA2Match_570/GetMatchHistory/V001/?league_id="

type SteamdotaRepository interface {
	GetLatestGameID() (string, error)
	GetAllMatches() (steamdota.Steamdota, error)
}

type steamdotaRepository struct {
	config     config
	logger     *zap.Logger
	httpClient http.Client
}

type config struct {
	token string
	id    string
}

func NewConfig(token string, id string) config {
	return config{
		token: token,
		id:    id,
	}
}

func NewSteamdotaRepository(logger *zap.Logger, httpClient http.Client, config config) SteamdotaRepository {
	return &steamdotaRepository{
		logger:     logger,
		httpClient: httpClient,
		config:     config,
	}
}

func (rx steamdotaRepository) GetLatestGameID() (string, error) {
	asd, _ := rx.GetAllMatches()

	return strconv.FormatInt(asd.Result.Matches[0].MatchID, 10), nil
}

func (rx steamdotaRepository) GetAllMatches() (steamdota.Steamdota, error) {
	s := fmt.Sprintf("%s%s&key=%s", endpoint, rx.config.id, rx.config.token)
	r, _ := rx.httpClient.Get(s)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.getAllPlayers failed")
		return steamdota.Steamdota{}, err
	}

	var o = steamdota.Steamdota{}
	if err := json.Unmarshal(body, &o); err != nil {
		rx.logger.Error("KungdotaRepository.getAllPlayers failed")
		return steamdota.Steamdota{}, err
	}

	return o, nil
}
