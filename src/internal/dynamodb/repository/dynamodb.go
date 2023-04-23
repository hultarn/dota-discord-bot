package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type DynamodbRepository interface {
	GetByCurrentWeekAndYear(ctx context.Context) (string, error)
	InsertCurrentWeekAndYear(ctx context.Context, id string) error
}

type dynamodbRepository struct {
	config config
	logger *zap.Logger
	client *dynamodb.Client
}

type config struct {
	c aws.Config
}

type Entry struct {
	Week         string `dynamodbav:"week"`
	Year         string `dynamodbav:"year"`
	MessageID    string `dynamodbav:"message_id"`
	CreationDate string `dynamodbav:"creation_date"`
}

func NewConfig(c aws.Config) config {
	return config{
		c: c,
	}
}

func NewDynamodbRepository(logger *zap.Logger, config config) DynamodbRepository {
	return &dynamodbRepository{
		config: config,
		logger: logger,
		client: dynamodb.NewFromConfig(config.c),
	}
}

func (rx dynamodbRepository) GetByCurrentWeekAndYear(ctx context.Context) (string, error) {
	var entry Entry

	year, week := time.Now().UTC().ISOWeek()

	item, err := rx.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String("discord"),
			Key: map[string]types.AttributeValue{
				"week": &types.AttributeValueMemberS{Value: strconv.Itoa(week)},
				"year": &types.AttributeValueMemberS{Value: strconv.Itoa(year)},
			},
		},
	)
	if err != nil {
		rx.logger.Error("DynamodbRepository.GetByCurrentWeekAndYear failed")
		return "", err
	}

	err = attributevalue.UnmarshalMap(item.Item, &entry)
	if err != nil {
		rx.logger.Error("DynamodbRepository.GetByCurrentWeekAndYear failed")
		return "", err
	}

	return entry.MessageID, nil
}

func (rx dynamodbRepository) InsertCurrentWeekAndYear(ctx context.Context, id string) error {
	t := time.Now()
	year, week := t.UTC().ISOWeek()

	e := &Entry{
		Week:         strconv.Itoa(week),
		Year:         strconv.Itoa(year),
		MessageID:    id,
		CreationDate: t.Format(time.RFC3339),
	}

	rx.logger.Info("DynamodbRepository.Insert")

	marshaledObjectToInsert, err := attributevalue.MarshalMap(e)
	if err != nil {
		rx.logger.Error("DynamodbRepository.GetByCurrentWeekAndYear failed")
		return err
	}

	_, err = rx.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String("discord"),
			Item:      marshaledObjectToInsert,
		},
	)
	if err != nil {
		rx.logger.Error("DynamodbRepository.GetByCurrentWeekAndYear failed")
		return err
	}

	return nil
}
