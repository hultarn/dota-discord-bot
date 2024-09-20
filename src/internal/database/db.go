package database

import (
	"context"
	"dota-discord-bot/src/internal/kungdota"

	"github.com/google/uuid"
)

type Message struct {
	MessageID    string    `dynamodbav:"message_id" json:"message_id" db:"message_id"`
	Week         string    `dynamodbav:"week" json:"week" db:"week"`
	Year         string    `dynamodbav:"year" json:"year" db:"year"`
	Game_1       uuid.UUID `dynamodbav:"game_1" json:"game_1" db:"game_1"`
	Game_2       uuid.UUID `dynamodbav:"game_2" json:"game_2" db:"game_2"`
	Game_3       uuid.UUID `dynamodbav:"game_3" json:"game_3" db:"game_3"`
	CreationDate string    `dynamodbav:"creation_date" json:"creation_date" db:"creation_date"`
}

type Sign struct {
	ID        int    `db:"id" json:"id"`
	GameID    string `db:"game_id" json:"game_id"`
	DiscordID string `db:"discord_id" json:"discord_id"`
}

type Repository interface {
	GetCurrent(ctx context.Context) (Message, error)
	GetByCurrentWeekAndYear(ctx context.Context) (string, error)
	InsertCurrentWeekAndYear(ctx context.Context, id string) error
	GetCurrentPlayers(ctx context.Context, idx int) ([]string, error)
	InsertPlayer(ctx context.Context, id string, i int) error
	ClearPlayers(ctx context.Context) error
	GetSigedPlayersByGame(ctx context.Context, id string) ([]Sign, error)
	GetShuffledTeams(ctx context.Context) (kungdota.ShuffledTeams, error)
	InsertShuffledPlayers(ctx context.Context, teams kungdota.ShuffledTeams, uuid uuid.UUID) error
}
