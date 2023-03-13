package repository

import (
	"context"
	"dota-discord-bot/src/internal/opendota"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type OpendotaRepository interface {
	GetMatch(ctx context.Context, id string) (opendota.OpenDotaGameObject, error)
	RequestMatch(id string) error
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

func (rx opendotaRepository) RequestMatch(id string) error {
	// TODO: Implement polling instead?
	_, err := rx.httpClient.Post(fmt.Sprintf("https://api.opendota.com/api/request/%s", id), "application/json", nil)
	if err != nil {
		rx.logger.Error(err.Error())

	}

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	rx.logger.Error(err.Error())

	// }

	// _, err = rx.httpClient.Get(fmt.Sprintf("https://api.opendota.com/api/request/%s", id))
	// if err != nil {
	// 	rx.logger.Error(err.Error())

	// }

	time.Sleep(time.Second * 60)

	return nil
}
