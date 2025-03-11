package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"skyhawk/backend/player/domain"
)

type Repository interface {
	SeasonStats(id string) (domain.PlayerSeasonStats, error)
	Save(ctx context.Context, tx *sql.Tx, player []domain.Player) (map[string]string, error)
}

type Repo struct {
	logger *zap.Logger
	db     *sqlx.DB
	redis  *redis.Client
}

const playerTtl = time.Hour * 24

func NewRepo(logger *zap.Logger, db *sqlx.DB, redis *redis.Client) Repository {

	return &Repo{logger: logger, db: db, redis: redis}
}

func (r *Repo) Save(ctx context.Context, tx *sql.Tx, players []domain.Player) (map[string]string, error) {
	// First check Redis for existing players in batch
	redisPipe := r.redis.Pipeline()
	redisResults := make(map[string]*redis.StringCmd)

	// Create player keys and check Redis first
	for i := range players {
		// Use composite key: "player:{name}:{team_id}"
		redisKey := fmt.Sprintf("player:%s:%s", players[i].Name, players[i].Team)
		redisResults[redisKey] = redisPipe.Get(ctx, redisKey)
	}

	// Execute all Redis gets in one batch
	_, err := redisPipe.Exec(ctx) // Fixed: removed * before err
	if err != nil && err != redis.Nil {
		r.logger.Warn("Redis pipeline error", zap.Error(err))
		// Continue execution even if Redis fails
	}

	// Process Redis results and collect players not in cache
	var missingPlayers []domain.Player
	playerIdsMap := make(map[string]string, len(players))

	for i := range players {
		redisKey := fmt.Sprintf("player:%s:%s", players[i].Name, players[i].Team)
		id, err := redisResults[redisKey].Result()
		if err == nil && id != "" {
			// Found in Redis
			players[i].ID = id
			playerIdsMap[players[i].Name] = id
		} else {
			// Not found in Redis, check DB before assigning new ID
			// Query for this specific player
			var playerID string
			err := tx.QueryRow("SELECT id FROM players WHERE name = ? AND team_id = ?",
				players[i].Name, players[i].Team).Scan(&playerID)
			if err == nil {
				// Found in DB
				players[i].ID = playerID
				playerIdsMap[players[i].Name] = playerID
				// Update Redis for this found player
				r.redis.Set(ctx, redisKey, playerID, 24*time.Hour)
			} else if err == sql.ErrNoRows {
				// Not found in DB either, assign new ID
				players[i].ID = uuid.New().String()
				missingPlayers = append(missingPlayers, players[i])
			} else {
				// Actual DB error
				return nil, fmt.Errorf("database error checking player: %w", err)
			}
		}
	}

	// Only insert players not found in Redis or in database
	if len(missingPlayers) > 0 {
		placeholders := make([]string, 0, len(missingPlayers))
		values := make([]interface{}, 0, len(missingPlayers)*3)

		for i := range missingPlayers {
			placeholders = append(placeholders, "(?, ?, ?)")
			values = append(values, missingPlayers[i].ID, missingPlayers[i].Name, missingPlayers[i].Team)
		}

		// Batch insert
		// Make sure team_id exists before inserting
		q := fmt.Sprintf("INSERT INTO players (id, name, team_id) VALUES %s ON DUPLICATE KEY UPDATE name=VALUES(name)",
			strings.Join(placeholders, ","))
		// Note: Ensure that the team_id values being inserted exist in the teams table to prevent foreign key errors
		_, err := tx.Exec(q, values...) // Fixed: removed * before err
		if err != nil {
			return nil, err
		}

		// Update Redis with new players
		updatePipe := r.redis.Pipeline()
		for i := range missingPlayers {
			redisKey := fmt.Sprintf("player:%s:%s", missingPlayers[i].Name, missingPlayers[i].Team)
			updatePipe.Set(ctx, redisKey, missingPlayers[i].ID, 24*time.Hour)
			playerIdsMap[missingPlayers[i].Name] = missingPlayers[i].ID
		}

		// Execute Redis updates
		if _, err := updatePipe.Exec(ctx); err != nil { // Fixed: removed * before err
			r.logger.Warn("Failed to update Redis with new players", zap.Error(err))
			// Continue even if Redis update fails
		}
	}

	return playerIdsMap, nil
}

func (r *Repo) SeasonStats(id string) (domain.PlayerSeasonStats, error) {
	var playerSeasonStatsDB PlayerSeasonStats
	row := r.db.QueryRow("select player_id, player_name, games_played, avg_points, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played  from player_season_stats where player_id = ?", id)

	if err := row.Scan(&playerSeasonStatsDB.PlayerID, &playerSeasonStatsDB.PlayerName, &playerSeasonStatsDB.GamesPlayed, &playerSeasonStatsDB.AvgPoints, &playerSeasonStatsDB.AvgRebounds, &playerSeasonStatsDB.AvgAssists, &playerSeasonStatsDB.AvgSteals, &playerSeasonStatsDB.AvgBlocks, &playerSeasonStatsDB.AvgFouls, &playerSeasonStatsDB.AvgTurnovers, &playerSeasonStatsDB.AvgMinutesPlayed); err != nil {
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
