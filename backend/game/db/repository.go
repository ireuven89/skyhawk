package db

import (
	"database/sql"
	"fmt"
	"strings"

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
	Save(tx *sql.Tx, game domain.GameStats) (string, error)
	Find(gameId string) (domain.GameStats, error)
	Begin() (*sql.Tx, error)
}

func NewRepo(db *sqlx.DB, logger *zap.Logger) Repo {

	return &Repository{db: db, logger: logger}
}

func (g *Repository) Begin() (*sql.Tx, error) {

	return g.db.Begin()
}

func (g *Repository) Save(tx *sql.Tx, game domain.GameStats) (string, error) {
	var placeHolders []string
	var values []interface{}

	if game.ID == "" {
		game.ID = uuid.New().String()
		game.GameID = uuid.New().String()
	}
	for _, team := range game.Teams {
		for _, player := range team.Players {
			placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			values = append(values, game.ID, player.ID, game.ID, game.Date, player.Points, player.Rebounds, player.Assists, player.Steals, player.Blocks, player.Fouls, player.Turnovers, player.MinutesPlayed)
		}
	}

	q := fmt.Sprintf("INSERT INTO game_stats (id, game_id, player_id, date, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played) values %s", strings.Join(placeHolders, ","))
	_, err := tx.Exec(q, values...)

	if err != nil {
		g.logger.Error("failed inserting stats", zap.Error(err))
		return "", err
	}

	return "", nil
}

func (g *Repository) Find(id string) (domain.GameStats, error) {
	var gameDB GameStatsDB

	row, err := g.db.Query("select id, game_id, player_id, date, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played from games where id = ?", id)

	if err != nil {
		return domain.GameStats{}, err
	}

	if err = row.Scan(&gameDB.ID, gameDB.Name); err != nil {
		return domain.GameStats{}, err
	}

	return domain.GameStats{}, nil
}
