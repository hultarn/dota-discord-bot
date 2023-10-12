package service

import (
	"context"
	"dota-discord-bot/src/internal/dynamodb/repository"

	"go.uber.org/zap"
)

type DynamodbService interface {
	GetCurrent(ctx context.Context) (repository.Entry, error)
	GetByCurrentWeekAndYear(ctx context.Context) (string, error)
	InsertCurrentWeekAndYear(ctx context.Context, id string) error
	GetCurrentPlayers(ctx context.Context, idx int) ([]string, error)
	InsertPlayer(ctx context.Context, id string, i int) error
	ClearPlayers(ctx context.Context) error
}

type dynamodbService struct {
	logger             *zap.Logger
	dynamodbRepository repository.DynamodbRepository
}

func NewDynamodbService(logger *zap.Logger, dynamodbRepository repository.DynamodbRepository) DynamodbService {
	return &dynamodbService{
		logger:             logger,
		dynamodbRepository: dynamodbRepository,
	}
}

func (rx dynamodbService) GetCurrent(ctx context.Context) (repository.Entry, error) {
	return rx.dynamodbRepository.GetCurrent(ctx)
}

func (rx dynamodbService) GetByCurrentWeekAndYear(ctx context.Context) (string, error) {
	return rx.dynamodbRepository.GetByCurrentWeekAndYear(ctx)
}

func (rx dynamodbService) InsertCurrentWeekAndYear(ctx context.Context, id string) error {
	return rx.dynamodbRepository.InsertCurrentWeekAndYear(ctx, id)
}

func (rx dynamodbService) GetCurrentPlayers(ctx context.Context, idx int) ([]string, error) {
	return rx.dynamodbRepository.GetCurrentPlayers(ctx, idx)
}

func (rx dynamodbService) InsertPlayer(ctx context.Context, id string, i int) error {
	return rx.dynamodbRepository.InsertPlayer(ctx, id, i)
}

func (rx dynamodbService) ClearPlayers(ctx context.Context) error {
	return rx.dynamodbRepository.ClearPlayers(ctx)
}
