package repository

import (
	"context"
	"database/sql"
	"dota-discord-bot/src/internal/database"
	"dota-discord-bot/src/internal/kungdota"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type mySqlRepository struct {
	config config
	logger *zap.Logger
	db     *sqlx.DB
}

type config struct {
	mySql string
}

func NewConfig(mySql string) config {
	return config{
		mySql: mySql,
	}
}

func NewMySqlRepository(logger *zap.Logger, config config) database.Repository {
	db, err := sqlx.Connect("mysql", config.mySql)
	if err != nil {
		log.Fatalln(err)
	}

	return &mySqlRepository{
		config: config,
		logger: logger,
		db:     db,
	}
}

func (rx mySqlRepository) GetCurrent(ctx context.Context) (database.Message, error) {
	year, week := time.Now().UTC().ISOWeek()

	var entry database.Message
	err := rx.db.Get(&entry, "SELECT * FROM messages WHERE year = ? AND week = ?", year, week)
	if err == sql.ErrNoRows {
		return database.Message{}, nil
	}

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("User: %+v\n", entry)

	return entry, nil
}

func (rx mySqlRepository) GetByCurrentWeekAndYear(ctx context.Context) (string, error) {
	entry, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("mySqlRepository.GetByCurrentWeekAndYear failed")
		return "", err
	}

	return entry.MessageID, nil
}

func (rx mySqlRepository) InsertCurrentWeekAndYear(ctx context.Context, id string) error {
	t := time.Now()
	year, week := t.UTC().ISOWeek()

	e := &database.Message{
		MessageID: id,
		Week:      strconv.Itoa(week),
		Year:      strconv.Itoa(year),
		Game_1:    uuid.New(),
		Game_2:    uuid.New(),
		Game_3:    uuid.New(),
	}

	rx.logger.Info("mySqlRepository.Insert")

	_, err := rx.db.Exec(
		"INSERT INTO messages (message_id, week, `year`, game_1, game_2, game_3) VALUES (?, ?, ?, ?, ?, ?);",
		e.MessageID, e.Week, e.Year, e.Game_1.String(), e.Game_2.String(), e.Game_3.String(),
	)
	if err != nil {
		rx.logger.Error("mySqlRepository.InsertCurrentWeekAndYear failed")
		return err
	}

	return nil
}

func (rx mySqlRepository) GetCurrentPlayers(ctx context.Context, idx int) ([]string, error) {
	entry, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("mySqlRepository.GetCurrentPlayers failed")
		return nil, err
	}

	var signs []database.Sign
	switch idx {
	case 1:
		signs, err = rx.GetSigedPlayersByGame(ctx, entry.Game_1.String())
	case 2:
		signs, err = rx.GetSigedPlayersByGame(ctx, entry.Game_2.String())
	case 3:
		signs, err = rx.GetSigedPlayersByGame(ctx, entry.Game_3.String())
	}
	if err != nil {
		rx.logger.Error("mySqlRepository.GetCurrentPlayers failed")
		return nil, err
	}

	var ret []string
	for _, sign := range signs {
		ret = append(ret, sign.DiscordID)
	}

	return ret, nil
}

func (rx mySqlRepository) InsertPlayer(ctx context.Context, id string, i int) error {
	e, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("InsertPlayer failed")
		return err
	}

	var uuid uuid.UUID
	switch i {
	case 0:
		uuid = e.Game_1
	case 1:
		uuid = e.Game_2
	case 2:
		uuid = e.Game_3
	default:
		rx.logger.Error("InsertPlayer failed")

		return fmt.Errorf("InsertPlayer invalid Id")
	}

	var sign database.Sign
	err = rx.db.Get(&sign, "SELECT * FROM signs WHERE game_id = ? AND discord_id = ?", uuid, id)
	if err == sql.ErrNoRows {
		_, err = rx.db.Exec(`INSERT INTO signs (game_id, discord_id) VALUES (?, ?);`, uuid, id)
		if err != nil {
			rx.logger.Error("InsertPlayer failed")
			return err
		}

		return nil
	}

	_, err = rx.db.Exec(`DELETE FROM signs WHERE game_id = ? AND discord_id = ?;`, uuid, id)
	if err != nil {
		rx.logger.Error("InsertPlayer failed")
		return err
	}

	return nil
}

func (rx mySqlRepository) ClearPlayers(ctx context.Context) error {
	e, err := rx.GetCurrent(ctx)
	if err != nil {
		rx.logger.Error("mySqlRepository.ClearPlayers failed")
		return err
	}

	rx.logger.Info("mySqlRepository.ClearPlayers")

	_, err = rx.db.Exec(`DELETE FROM signs WHERE game_id = ? AND discord_id = ?;`, e, e)
	if err != nil {
		rx.logger.Error("InsertCurrentWeekAndYear failed")
		return err
	}

	return nil
}

func (rx mySqlRepository) GetSigedPlayersByGame(ctx context.Context, id string) ([]database.Sign, error) {
	rows, err := rx.db.Query("SELECT * FROM signs WHERE game_id = ?", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var entries []database.Sign
	for rows.Next() {
		var sign database.Sign
		if err := rows.Scan(&sign.ID, &sign.GameID, &sign.DiscordID); err != nil {
			panic(err)
		}
		entries = append(entries, sign)
	}

	return entries, nil
}

func (rx mySqlRepository) GetShuffledTeams(ctx context.Context) (kungdota.ShuffledTeams, error) {
	type shuffled struct {
		team      int
		discordID string
	}

	rows, err := rx.db.Query("SELECT team, discord_id FROM shuffled_teams ORDER BY creation_date DESC LIMIT 10;")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	shuffledTeams := kungdota.ShuffledTeams{}
	for rows.Next() {
		var player shuffled
		if err := rows.Scan(&player.team, &player.discordID); err != nil {
			panic(err)
		}

		if player.team == 0 {
			shuffledTeams.TeamOne.Players = append(shuffledTeams.TeamOne.Players, kungdota.Player{
				DiscordID: player.discordID,
			})
		} else {
			shuffledTeams.TeamTwo.Players = append(shuffledTeams.TeamTwo.Players, kungdota.Player{
				DiscordID: player.discordID,
			})
		}
	}

	return shuffledTeams, nil
}

func (rx mySqlRepository) InsertShuffledPlayers(ctx context.Context, teams kungdota.ShuffledTeams, uuid uuid.UUID) error {
	for _, teamOne := range teams.TeamOne.Players {
		_, err := rx.db.Exec(`
        INSERT INTO shuffled_teams (shuffle_id, team, discord_id)
        VALUES (?, ?, ?)`, uuid, 0, teamOne.DiscordID)
		if err != nil {
			return fmt.Errorf("failed to insert player: %v", err)
		}

	}

	for _, teamTwo := range teams.TeamTwo.Players {
		_, err := rx.db.Exec(`
        INSERT INTO shuffled_teams (shuffle_id, team, discord_id)
        VALUES (?, ?, ?)`, uuid, 1, teamTwo.DiscordID)
		if err != nil {
			return fmt.Errorf("failed to insert player: %v", err)
		}
	}

	return nil
}
