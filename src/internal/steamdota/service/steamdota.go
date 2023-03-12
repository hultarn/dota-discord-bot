package service

import (
	"dota-discord-bot/src/internal/steamdota/repository"

	"go.uber.org/zap"
)

type SteamdotaService interface {
	GetLatestGameID() (string, error)
	// GetProperties() Properties
}

type steamdotaService struct {
	logger              *zap.Logger
	steamdotaRepository repository.SteamdotaRepository
	// properties         Properties
}

//priv?
// type Properties struct {
// 	ShuffledTeams kungdota.ShuffledTeams
// }

func NewSteamdotaService(logger *zap.Logger, steamdotaRepository repository.SteamdotaRepository) SteamdotaService {
	return &steamdotaService{
		logger:              logger,
		steamdotaRepository: steamdotaRepository,
	}
}

func (rx steamdotaService) GetLatestGameID() (string, error) {
	r, _ := rx.steamdotaRepository.GetLatestGameID()

	return r, nil
}

// func (rx *steamdotaService) GetProperties() Properties {
// 	return rx.properties
// }
