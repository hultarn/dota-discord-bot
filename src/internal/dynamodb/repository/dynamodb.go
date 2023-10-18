package repository

import (
	"context"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type DynamodbRepository interface {
	GetCurrent(ctx context.Context) (Entry, error)
	GetByCurrentWeekAndYear(ctx context.Context) (string, error)
	InsertCurrentWeekAndYear(ctx context.Context, id string) error
	GetCurrentPlayers(ctx context.Context, idx int) ([]string, error)
	InsertPlayer(ctx context.Context, id string, i int) error
	ClearPlayers(ctx context.Context) error
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
	Week         string   `dynamodbav:"week"`
	Year         string   `dynamodbav:"year"`
	MessageID    string   `dynamodbav:"message_id"`
	CreationDate string   `dynamodbav:"creation_date"`
	Game_1       []string `dynamodbav:"game_1"`
	Game_2       []string `dynamodbav:"game_2"`
	Game_3       []string `dynamodbav:"game_3"`
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

func (rx dynamodbRepository) GetCurrent(ctx context.Context) (Entry, error) {
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
		rx.logger.Error("DynamodbRepository.GetCurrent.GetItem failed")
		return Entry{}, err
	}

	err = attributevalue.UnmarshalMap(item.Item, &entry)
	if err != nil {
		rx.logger.Error("DynamodbRepository.GetCurrent.UnmarshalMap failed")
		return Entry{}, err
	}

	return entry, nil
}

func (rx dynamodbRepository) GetByCurrentWeekAndYear(ctx context.Context) (string, error) {
	entry, err := rx.GetCurrent(ctx)
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
		rx.logger.Error("DynamodbRepository.InsertCurrentWeekAndYear failed")
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
		rx.logger.Error("DynamodbRepository.InsertCurrentWeekAndYear failed")
		return err
	}

	return nil
}

func (rx dynamodbRepository) GetCurrentPlayers(ctx context.Context, idx int) ([]string, error) {
	entry, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("DynamodbRepository.GetCurrentPlayers failed")
		return nil, err
	}

	var ret []string

	switch idx {
	case 1:
		ret = entry.Game_1
	case 2:
		ret = entry.Game_2
	case 3:
		ret = entry.Game_3
	}

	return ret, nil
}

func (rx dynamodbRepository) InsertPlayer(ctx context.Context, id string, i int) error {
	e, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("DynamodbRepository.InsertPlayer failed")
		return err
	}

	var tmp *[]string
	switch i {
	case 0:
		tmp = &e.Game_1
	case 1:
		tmp = &e.Game_2
	case 2:
		tmp = &e.Game_3
	default:
		rx.logger.Error("DynamodbRepository.InsertPlayer failed")
		return err
	}

	if slices.Contains(*tmp, id) {
		rx.logger.Error("DynamodbRepository.InsertPlayer failed, player already signed up")
		return nil
	}

	if len(*tmp) >= 10 {
		rx.logger.Error("DynamodbRepository.InsertPlayer failed, game is full")
		return nil
	}

	*tmp = append(*tmp, id)

	rx.logger.Info("DynamodbRepository.Insert")

	marshaledObjectToInsert, err := attributevalue.MarshalMap(e)
	if err != nil {
		rx.logger.Error("DynamodbRepository.InsertPlayer failed")
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
		rx.logger.Error("DynamodbRepository.InsertPlayer failed")
		return err
	}

	return nil
}

func (rx dynamodbRepository) ClearPlayers(ctx context.Context) error {
	e, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("DynamodbRepository.ClearPlayers failed")
		return err
	}

	rx.logger.Info("DynamodbRepository.ClearPlayers")

	e.Game_1 = []string{}
	e.Game_2 = []string{}
	e.Game_3 = []string{}

	marshaledObjectToInsert, err := attributevalue.MarshalMap(e)
	if err != nil {
		rx.logger.Error("DynamodbRepository.InsertPlayer failed")
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
		rx.logger.Error("DynamodbRepository.ClearPlayers failed")
		return err
	}

	return nil
}
