package service

import (
	"context"
	"dota-discord-bot/src/internal/kungdota"
	"dota-discord-bot/src/internal/opendota"
	"dota-discord-bot/src/internal/opendota/repository"

	"go.uber.org/zap"
)

type OpendotaService interface {
	GetMatch(id string) (opendota.OpenDotaGameObject, error)
	RequestMatch(id string) error
}

type opendotaService struct {
	logger             *zap.Logger
	opendotaRepository repository.OpendotaRepository
}

type Properties struct {
	ShuffledTeams kungdota.ShuffledTeams
}

func NewOpendotaService(logger *zap.Logger, opendotaRepository repository.OpendotaRepository) OpendotaService {
	return &opendotaService{
		logger:             logger,
		opendotaRepository: opendotaRepository,
	}
}

func (rx opendotaService) GetMatch(id string) (opendota.OpenDotaGameObject, error) {
	return rx.opendotaRepository.GetMatch(context.Background(), id)
}

func (rx opendotaService) RequestMatch(id string) error {
	return rx.opendotaRepository.RequestMatch(id)
}
