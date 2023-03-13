package repository

import (
	"bytes"
	"context"
	"dota-discord-bot/src/internal/kungdota"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type KungdotaRepository interface {
	GetByNames(ctx context.Context, ids []string) (kungdota.Players2, error)
	GetByDiscordID(ctx context.Context, ids []string) ([]kungdota.Players2, error)
	PostMatch(m kungdota.Match) error
	GetAllPlayers() (kungdota.Players2, error)
}

type kungdotaRepository struct {
	config     config
	logger     *zap.Logger
	httpClient http.Client
}

type config struct {
	leagueID string
}

func NewConfig(leagueID string) config {
	return config{
		leagueID: leagueID,
	}
}

func NewKungdotaRepository(logger *zap.Logger, httpClient http.Client, config config) KungdotaRepository {
	return &kungdotaRepository{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

func (rx kungdotaRepository) PostMatch(m kungdota.Match) error {
	asd, _ := json.Marshal(m)
	fmt.Println(m)
	buff := bytes.NewBuffer(asd)
	fmt.Println(buff)
	tmp, err := rx.httpClient.Post("https://api.bollsvenskan.jacobadlers.com/match/", "application/json", buff)
	fmt.Println(tmp)

	return err
}

func (rx kungdotaRepository) GetAllPlayers() (kungdota.Players2, error) {
	resp, _ := rx.httpClient.Get("https://api.bollsvenskan.jacobadlers.com/player")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.GetAllPlayers failed")
		return kungdota.Players2{}, err
	}

	var players = kungdota.Players2{}
	if err := json.Unmarshal(body, &players); err != nil {
		rx.logger.Error("KungdotaRepository.GetAllPlayers failed")
		return kungdota.Players2{}, err
	}

	return players, nil
}

func (rx kungdotaRepository) GetByNames(ctx context.Context, ids []string) (kungdota.Players2, error) {
	players, err := rx.GetAllPlayers()
	if err != nil {
		rx.logger.Error("KungdotaRepository.GetByNames failed")
		return kungdota.Players2{}, err
	}

	pList := make([]kungdota.Players, 0)
	for _, id := range ids {
		for _, p := range players.Players {
			if p.Username == id {
				pList = append(pList, p)
				goto found
			}
		}
		rx.logger.Error("KungdotaRepository.GetByNames player not found")
		return kungdota.Players2{}, errors.New("player not found")
	found:
	}

	return kungdota.Players2{
		Players: pList,
	}, nil
}

func (rx kungdotaRepository) GetByDiscordID(ctx context.Context, ids []string) ([]kungdota.Players2, error) {
	return nil, nil
}
