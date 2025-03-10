package db

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"skyhawk/backend/player/domain"
)

type Repository interface {
	SeasonStats(id string) (domain.PlayerSeasonStats, error)
	Save(tx *sql.Tx, player domain.Player) (string, error)
}

type Repo struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewRepo(logger *zap.Logger, db *sqlx.DB) Repository {

	return &Repo{logger: logger, db: db}
}

func (r *Repo) Save(tx *sql.Tx, player domain.Player) (string, error) {
	if player.ID == "" {
		player.ID = uuid.New().String()
	}
	res, err := tx.Query("INSERT INTO players (id, name, team_id) VALUES (?,?, ?) ON DUPLICATE KEY UPDATE name = ?", player.ID, player.Name, player.Team, player.Name)

	if err != nil {
		r.logger.Error("failed inserting rows", zap.Error(err))
		return "", err
	}

	if res.Err() != nil {
		return "", err
	}

	return player.ID, nil
}

func (r *Repo) SeasonStats(id string) (domain.PlayerSeasonStats, error) {
	var playerSeasonStatsDB PlayerSeasonStats
	row := r.db.QueryRow("select player_id, player_name, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played  from player_season_stats where player_id = ?", id)

	if err := row.Scan(&playerSeasonStatsDB.PlayerID, &playerSeasonStatsDB.PlayerName, &playerSeasonStatsDB.AvgRebounds, &playerSeasonStatsDB.AvgAssists, &playerSeasonStatsDB.AvgSteals, &playerSeasonStatsDB.AvgBlocks, &playerSeasonStatsDB.AvgFouls, &playerSeasonStatsDB.AvgTurnovers, &playerSeasonStatsDB.AvgMinutesPlayed); err != nil {
		return domain.PlayerSeasonStats{}, err
	}

	return toDomain(playerSeasonStatsDB), nil
}

func toDomain(dbModel PlayerSeasonStats) domain.PlayerSeasonStats {

	return domain.PlayerSeasonStats{
		PlayerID:         dbModel.PlayerID,
		PlayerName:       dbModel.PlayerName,
		TeamID:           dbModel.TeamID,
		TeamName:         dbModel.TeamName,
		GamesPlayed:      dbModel.GamesPlayed,
		AvgPoints:        dbModel.AvgPoints,
		AvgRebounds:      dbModel.AvgRebounds,
		AvgAssists:       dbModel.AvgAssists,
		AvgSteals:        dbModel.AvgSteals,
		AvgBlocks:        dbModel.AvgBlocks,
		AvgFouls:         dbModel.AvgFouls,
		AvgTurnovers:     dbModel.AvgTurnovers,
		AvgMinutesPlayed: dbModel.AvgMinutesPlayed,
	}
}
