package repository

import (
	"context"
	"dota-discord-bot/src/internal/opendota"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type OpendotaRepository interface {
	GetMatch(ctx context.Context, id string) (opendota.OpenDotaGameObject, error)
}

type opendotaRepository struct {
	config     config
	logger     *zap.Logger
	httpClient http.Client
}

type config struct {
	token string
}

func NewConfig(token string) config {
	return config{
		token: token,
	}
}

func NewOpendotaRepository(logger *zap.Logger, httpClient http.Client, config config) OpendotaRepository {
	return &opendotaRepository{
		logger:     logger,
		httpClient: httpClient,
		config:     config,
	}
}

func (rx opendotaRepository) GetMatch(ctx context.Context, id string) (opendota.OpenDotaGameObject, error) {
	resp, err := rx.httpClient.Get(fmt.Sprintf("https://api.opendota.com/api/matches/%s?key=%s", id, rx.config.token))
	if err != nil {
		rx.logger.Error(err.Error())
		return opendota.OpenDotaGameObject{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rx.logger.Error(err.Error())
		return opendota.OpenDotaGameObject{}, err
	}

	var r = &opendota.OpenDotaGameObject{}
	json.Unmarshal(body, &r)

	return *r, nil
}
