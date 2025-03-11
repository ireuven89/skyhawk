package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"skyhawk/backend/team/domain"
	"testing"
)

func TestRepo_New(t *testing.T) {
	// Setup
	db, _, _ := createMockDB(t)
	rdb := createMockRedis(t)
	logger := zaptest.NewLogger(t)

	// Test Redis connection is working
	ctx := context.Background()
	require.NoError(t, rdb.Ping(ctx).Err(), "Redis ping should succeed")

	// Test
	repo := New(db, rdb, logger)

	// Assert
	assert.NotNil(t, repo, "Repository should not be nil")
}

func TestRepo_Save(t *testing.T) {
	t.Run("team found in redis", func(t *testing.T) {
		// Setup
		db, _, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx := context.Background()
		teamName := "Lakers"
		teamID := uuid.New().String()

		// Set redis mock to return the team ID
		require.NoError(t, rdb.Set(ctx, teamName, teamID, 0).Err())

		// Test
		id, err := repo.Save(ctx, nil, domain.Team{Name: teamName})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, teamID, id)

		// Verify Redis was used (this GET is just for verification)
		val, err := rdb.Get(ctx, teamName).Result()
		assert.NoError(t, err)
		assert.Equal(t, teamID, val)
	})

	t.Run("redis error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)

		// Create a controlled redis mock that returns an error
		mr := miniredis.NewMiniRedis()
		require.NoError(t, mr.Start())
		defer mr.Close()
		mr.SetError("forced error")

		rdb := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})

		// Initialize Redis connection (even though it will error)
		ctx := context.Background()
		_ = rdb.Ping(ctx)

		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx = context.Background()
		teamName := "Lakers"
		teamID := uuid.New().String()

		// DB will be checked since Redis failed
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(teamID, teamName)

		dbMock.ExpectQuery("SELECT id, name FROM teams WHERE name = ?").
			WithArgs(teamName).
			WillReturnRows(rows)

		// Test
		id, err := repo.Save(ctx, nil, domain.Team{Name: teamName})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, teamID, id)

		// After finding in DB, code should try to update Redis cache (which will fail due to our forced error)
		// Because this will fail silently, we just confirm the result was still returned properly
	})

	t.Run("team found in DB", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx := context.Background()
		teamName := "Lakers"
		teamID := uuid.New().String()

		// Check that Redis doesn't have the team initially
		_, err := rdb.Get(ctx, teamName).Result()
		assert.Equal(t, redis.Nil, err, "Redis should not have team before test")

		// DB mock will return existing team
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(teamID, teamName)

		dbMock.ExpectQuery("SELECT id, name FROM teams WHERE name = ?").
			WithArgs(teamName).
			WillReturnRows(rows)

		// Test
		id, err := repo.Save(ctx, nil, domain.Team{Name: teamName})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, teamID, id)

		// Verify Redis has been updated after DB lookup
		val, err := rdb.Get(ctx, teamName).Result()
		assert.NoError(t, err)
		assert.Equal(t, teamID, val, "Redis should be updated with team ID from DB")
	})

	t.Run("team not found, creating new", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx := context.Background()
		// Mock the transaction
		dbMock.ExpectBegin()
		mockTx, err := db.Begin()
		require.NoError(t, err)

		teamName := "New Team"

		// DB mock will return no existing team
		rows := sqlmock.NewRows([]string{"id", "name"})

		dbMock.ExpectQuery("SELECT id, name FROM teams WHERE name = ?").
			WithArgs(teamName).
			WillReturnRows(rows)

		// Expect INSERT
		dbMock.ExpectExec("INSERT INTO teams \\(id, name\\) VALUES \\(\\?,\\?\\)").
			WithArgs(sqlmock.AnyArg(), teamName).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Test
		id, err := repo.Save(ctx, mockTx, domain.Team{Name: teamName})

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// Verify Redis has been updated
		val, err := rdb.Get(ctx, teamName).Result()
		assert.NoError(t, err)
		assert.Equal(t, id, val)
	})

	t.Run("DB query error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx := context.Background()
		teamName := "Error Team"

		// DB mock will return error
		dbMock.ExpectQuery("SELECT id, name FROM teams WHERE name = ?").
			WithArgs(teamName).
			WillReturnError(errors.New("db query error"))

		// Test
		id, err := repo.Save(ctx, nil, domain.Team{Name: teamName})

		// Assert
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("insert error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx := context.Background()
		// Mock the transaction
		dbMock.ExpectBegin()
		mockTx, err := db.Begin()
		require.NoError(t, err)

		teamName := "Error Team"

		// DB mock will return no existing team
		rows := sqlmock.NewRows([]string{"id", "name"})

		dbMock.ExpectQuery("SELECT id, name FROM teams WHERE name = ?").
			WithArgs(teamName).
			WillReturnRows(rows)

		// Expect INSERT to fail
		dbMock.ExpectExec("INSERT INTO teams \\(id, name\\) VALUES \\(\\?,\\?\\)").
			WithArgs(sqlmock.AnyArg(), teamName).
			WillReturnError(errors.New("insert error"))

		// Test
		id, err := repo.Save(ctx, mockTx, domain.Team{Name: teamName})

		// Assert
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("scan error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		ctx := context.Background()
		teamName := "Invalid Team"

		// Check that Redis doesn't have the team initially
		_, err := rdb.Get(ctx, teamName).Result()
		assert.Equal(t, redis.Nil, err, "Redis should not have team before test")

		// Return invalid column type to cause scan error
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nil, teamName) // nil ID will cause scan error

		dbMock.ExpectQuery("SELECT id, name FROM teams WHERE name = ?").
			WithArgs(teamName).
			WillReturnRows(rows)

		// Test
		id, err := repo.Save(ctx, nil, domain.Team{Name: teamName})

		// Assert
		assert.Error(t, err)
		assert.Empty(t, id)

		// Verify Redis still doesn't have the value
		_, err = rdb.Get(ctx, teamName).Result()
		assert.Equal(t, redis.Nil, err, "Redis should still not have team after error")
	})
}

func TestRepo_Find(t *testing.T) {
	t.Run("team found", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		teamID := uuid.New().String()
		teamName := "Lakers"

		// DB mock will return team
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(teamID, teamName)

		dbMock.ExpectQuery("select id, name from teams where name = ?").
			WithArgs(teamID).
			WillReturnRows(rows)

		// Test
		team, err := repo.Find(teamID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, teamID, team.ID)
		assert.Equal(t, teamName, team.Name)
	})

	t.Run("team not found", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		teamID := "nonexistent"

		// DB mock will return no team
		dbMock.ExpectQuery("select id, name from teams where name = ?").
			WithArgs(teamID).
			WillReturnError(sql.ErrNoRows)

		// Test
		team, err := repo.Find(teamID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.Team{}, team)
	})
}

func TestRepo_GetStats(t *testing.T) {
	t.Run("stats found", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		teamID := uuid.New().String()
		teamName := "Lakers"

		// Mock stats data
		stats := SeasonStats{
			TeamID:           teamID,
			TeamName:         teamName,
			GamesPlayed:      82,
			AvgPoints:        105.6,
			AvgRebounds:      42.3,
			AvgAssists:       24.5,
			AvgSteals:        8.2,
			AvgBlocks:        5.1,
			AvgFouls:         19.4,
			AvgTurnovers:     13.2,
			AvgMinutesPlayed: 240.0,
		}

		// DB mock will return stats - match the SQL and fields in the GetStats method
		rows := sqlmock.NewRows([]string{
			"team_id", "team_name", "avg_rebounds", "avg_assists",
			"avg_steals", "avg_blocks", "avg_fouls", "avg_turnovers", "avg_minutes_played",
		}).AddRow(
			stats.TeamID, stats.TeamName, stats.AvgRebounds,
			stats.AvgAssists, stats.AvgSteals, stats.AvgBlocks, stats.AvgFouls, stats.AvgTurnovers,
			stats.AvgMinutesPlayed,
		)

		dbMock.ExpectQuery("select team_id, team_name, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played from player_season_stats where team_id = ?").
			WithArgs(teamID).
			WillReturnRows(rows)

		// Test
		result, err := repo.GetStats(teamID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, teamID, result.TeamID)
		assert.Equal(t, teamName, result.TeamName)
		assert.Equal(t, 42.3, result.AvgRebounds)
		assert.Equal(t, 24.5, result.AvgAssists)
		assert.Equal(t, 8.2, result.AvgSteals)
		assert.Equal(t, 5.1, result.AvgBlocks)
		assert.Equal(t, 19.4, result.AvgFouls)
		assert.Equal(t, 13.2, result.AvgTurnovers)
		assert.Equal(t, 240.0, result.AvgMinutesPlayed)
	})

	t.Run("stats not found", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		rdb := createMockRedis(t)
		logger := zaptest.NewLogger(t)
		repo := New(db, rdb, logger)

		teamID := "nonexistent"

		// DB mock will return no stats - use the exact SQL query
		dbMock.ExpectQuery("select team_id, team_name, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played from player_season_stats where team_id = ?").
			WithArgs(teamID).
			WillReturnError(sql.ErrNoRows)

		// Test
		stats, err := repo.GetStats(teamID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.SeasonStats{}, stats)
	})
}

// Helper functions for creating mocks
func createMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, err
}

func createMockRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create mock Redis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Initialize Redis connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to ping Redis: %v", err)
	}

	// Clean up when test is done
	t.Cleanup(func() {
		client.Close()
		mr.Close()
	})

	return client
}
