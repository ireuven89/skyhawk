package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"skyhawk/backend/game/domain"
)

type Repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

type Repo interface {
	Save(tx *sql.Tx, game domain.GameStatsReq) (string, error)
	Find(gameId string) ([]domain.GameStats, error)
	Begin() (*sql.Tx, error)
}

func NewRepo(db *sqlx.DB, logger *zap.Logger) Repo {

	return &Repository{db: db, logger: logger}
}

func (g *Repository) Begin() (*sql.Tx, error) {

	return g.db.Begin()
}

func (g *Repository) Save(tx *sql.Tx, game domain.GameStatsReq) (string, error) {
	var placeHolders []string
	var values []interface{}
	gameId := uuid.New().String()

	// Check if we have any players to insert before proceeding
	playerCount := 0
	for _, team := range game.Teams {
		playerCount += len(team.Players)
	}

	if playerCount == 0 {
		return gameId, nil
	}

	for _, team := range game.Teams {
		for _, player := range team.Players {
			id := uuid.New().String()
			placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			values = append(values, id, gameId, player.ID, game.Date, player.Points, player.Rebounds, player.Assists, player.Steals, player.Blocks, player.Fouls, player.Turnovers, player.MinutesPlayed)
		}
	}

	q := fmt.Sprintf("INSERT INTO game_stats (id, game_id, player_id, date, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played) values %s", strings.Join(placeHolders, ","))
	_, err := tx.Exec(q, values...)
	if err != nil {
		g.logger.Error("failed inserting stats", zap.Error(err))
		return "", err
	}
	return gameId, nil
}

func (g *Repository) Find(id string) ([]domain.GameStats, error) {
	var result []domain.GameStats
	// Fix: Changed game*id to game_id
	row, err := g.db.Query("select g.game_id, g.player_id, p.name, g.date, g.points, g.rebounds, g.assists, g.steals, g.blocks, g.fouls, g.turnovers, g.minutes_played from game_stats g join players p on player_id = p.id where game_id =? ", id)
	if err != nil {
		return nil, err
	}
	defer row.Close() // Add this to ensure resources are properly released

	for row.Next() {
		gameDB := GameStatsDB{}
		if err = row.Scan(&gameDB.ID, &gameDB.PlayerID, &gameDB.Name, &gameDB.Date, &gameDB.Points, &gameDB.Rebounds, &gameDB.Assists, &gameDB.Steals, &gameDB.Blocks, &gameDB.Fouls, &gameDB.Turnovers, &gameDB.MinutesPlayed); err != nil {
			return nil, err
		}
		game, err := toDomain(gameDB)
		if err != nil {
			return nil, err
		}
		result = append(result, game)
	}
	return result, nil
}

func toDomain(db GameStatsDB) (domain.GameStats, error) {
	parsedDate, err := time.Parse(time.DateTime, db.Date)

	if err != nil {
		return domain.GameStats{}, err
	}

	return domain.GameStats{
		ID:            db.ID,
		Date:          parsedDate,
		Name:          db.Name,
		PlayerID:      db.PlayerID,
		Points:        db.Points,
		Steals:        db.Steals,
		Fouls:         db.Fouls,
		Turnovers:     db.Turnovers,
		Blocks:        db.Blocks,
		Assists:       db.Assists,
		MinutesPlayed: db.MinutesPlayed,
		Rebounds:      db.Rebounds,
	}, nil
}
