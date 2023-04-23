package service

import (
	"context"
	"dota-discord-bot/src/internal/dynamodb/repository"

	"go.uber.org/zap"
)

type DynamodbService interface {
	GetByCurrentWeekAndYear(ctx context.Context) (string, error)
	InsertCurrentWeekAndYear(ctx context.Context, id string) error
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

func (rx dynamodbService) GetByCurrentWeekAndYear(ctx context.Context) (string, error) {
	return rx.dynamodbRepository.GetByCurrentWeekAndYear(ctx)
}

func (rx dynamodbService) InsertCurrentWeekAndYear(ctx context.Context, id string) error {
	return rx.dynamodbRepository.InsertCurrentWeekAndYear(ctx, id)
}
