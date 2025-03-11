package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"skyhawk/backend/team/domain"
)

const timeTtl = time.Minute * 5

type Repository interface {
	Save(ctx context.Context, tx *sql.Tx, team domain.Team) (string, error)
	Find(id string) (domain.Team, error)
	GetStats(id string) (domain.SeasonStats, error)
}

type Repo struct {
	db     *sqlx.DB
	logger *zap.Logger
	redis  *redis.Client
}

func New(db *sqlx.DB, redis *redis.Client, logger *zap.Logger) Repository {

	return &Repo{
		db:     db,
		logger: logger,
		redis:  redis,
	}
}

func (r *Repo) Save(ctx context.Context, tx *sql.Tx, team domain.Team) (string, error) {

	//check if exists in redis - to reduce lattency and db overload
	// Check if exists in Redis
	res, err := r.redis.Get(ctx, team.Name).Result()
	// Key exists in Redis
	if err == nil {
		return res, nil
	}

	if err != redis.Nil {
		// Unexpected Redis error
		r.logger.Warn("Redis error", zap.Error(err), zap.String("team", team.Name))
	}

	r.logger.Info("Redis miss, checking DB", zap.String("team", team.Name))

	// Check if exists in DB
	row, err := r.db.Query("SELECT id, name FROM teams WHERE name = ?", team.Name)
	if err != nil {
		return "", err
	}
	defer row.Close() // Ensure rows are closed

	if !row.Next() {
		// Team doesn't exist, create new
		id := uuid.New().String()
		_, err = tx.Exec("INSERT INTO teams (id, name) VALUES (?,?)", id, team.Name)
		if err != nil {
			r.logger.Error("Failed inserting team", zap.Error(err))
			return "", err
		}

		// Set in Redis
		err = r.redis.Set(ctx, team.Name, id, timeTtl).Err()
		if err != nil {
			r.logger.Warn("Failed inserting to Redis", zap.Error(err))
		}

		return id, nil
	}

	// Team exists, get ID from DB
	var teamDB Team
	if err = row.Scan(&teamDB.ID, &teamDB.Name); err != nil {
		return "", err
	}

	// Set in Redis
	err = r.redis.Set(ctx, team.Name, teamDB.ID, timeTtl).Err()
	if err != nil {
		r.logger.Warn("Failed inserting to Redis", zap.Error(err))
	}

	return teamDB.ID, nil
}

func (r *Repo) Find(id string) (domain.Team, error) {
	var teamDB Team

	row := r.db.QueryRow("select id, name  from teams where name = ?", id)

	if err := row.Scan(&teamDB.ID, &teamDB.Name); err != nil {
		return domain.Team{}, err
	}

	return domain.Team{
		ID:   teamDB.ID,
		Name: teamDB.Name,
	}, nil
}

func (r *Repo) GetStats(id string) (domain.SeasonStats, error) {
	var seasonDB SeasonStats

	row := r.db.QueryRow("select team_id, team_name, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played  from player_season_stats where team_id = ?", id)

	if err := row.Scan(
		&seasonDB.TeamID,
		&seasonDB.TeamName,
		&seasonDB.AvgRebounds,
		&seasonDB.AvgAssists,
		&seasonDB.AvgSteals,
		&seasonDB.AvgBlocks,
		&seasonDB.AvgFouls,
		&seasonDB.AvgTurnovers,
		&seasonDB.AvgMinutesPlayed,
	); err != nil {
		return domain.SeasonStats{}, err
	}

	return toDomain(seasonDB), nil
}

func toDomain(team SeasonStats) domain.SeasonStats {

	return domain.SeasonStats{
		TeamID:           team.TeamID,
		TeamName:         team.TeamName,
		GamesPlayed:      team.GamesPlayed,
		AvgPoints:        team.AvgPoints,
		AvgRebounds:      team.AvgRebounds,
		AvgAssists:       team.AvgAssists,
		AvgSteals:        team.AvgSteals,
		AvgBlocks:        team.AvgBlocks,
		AvgFouls:         team.AvgFouls,
		AvgTurnovers:     team.AvgTurnovers,
		AvgMinutesPlayed: team.AvgMinutesPlayed,
	}
}
