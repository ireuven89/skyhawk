package db

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"skyhawk/backend/game/domain"
)

func TestRepository_Begin(t *testing.T) {
	// Setup
	db, dbMock, _ := createMockDB(t)
	logger := zaptest.NewLogger(t)
	repo := NewRepo(db, logger)

	// Set up mock to expect a transaction
	dbMock.ExpectBegin()

	// Test
	tx, err := repo.Begin()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestRepository_Save(t *testing.T) {
	t.Run("successful save", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		// Create transaction mock
		dbMock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		// Prepare test data
		gameDate := time.Now()
		gameReq := domain.GameStatsReq{
			Date: gameDate,
			Teams: []domain.Team{
				{
					ID: "team1",
					Players: []domain.Player{
						{
							ID:            "player1",
							Points:        20,
							Rebounds:      10,
							Assists:       5,
							Steals:        2,
							Blocks:        1,
							Fouls:         3,
							Turnovers:     2,
							MinutesPlayed: 35,
						},
						{
							ID:            "player2",
							Points:        15,
							Rebounds:      8,
							Assists:       7,
							Steals:        1,
							Blocks:        0,
							Fouls:         2,
							Turnovers:     1,
							MinutesPlayed: 30,
						},
					},
				},
				{
					ID: "team2",
					Players: []domain.Player{
						{
							ID:            "player3",
							Points:        18,
							Rebounds:      9,
							Assists:       4,
							Steals:        3,
							Blocks:        2,
							Fouls:         2,
							Turnovers:     3,
							MinutesPlayed: 32,
						},
					},
				},
			},
		}

		// Set up expectations for the INSERT statement
		// Using a generic regex for the INSERT statement to avoid strict matching
		dbMock.ExpectExec("INSERT INTO game_stats").
			WillReturnResult(sqlmock.NewResult(1, 3)) // 3 rows affected (3 players)

		// Test
		gameID, err := repo.Save(tx, gameReq)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, gameID)
	})

	t.Run("insert error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		// Create transaction mock
		dbMock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		// Prepare test data
		gameDate := time.Now()
		gameReq := domain.GameStatsReq{
			Date: gameDate,
			Teams: []domain.Team{
				{
					ID: "team1",
					Players: []domain.Player{
						{
							ID:            "player1",
							Points:        20,
							Rebounds:      10,
							Assists:       5,
							Steals:        2,
							Blocks:        1,
							Fouls:         3,
							Turnovers:     2,
							MinutesPlayed: 35,
						},
					},
				},
			},
		}

		// Set up expectations for the INSERT statement to fail
		dbMock.ExpectExec("INSERT INTO game_stats").
			WillReturnError(sql.ErrConnDone)

		// Test
		gameID, err := repo.Save(tx, gameReq)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, gameID)
	})

	t.Run("empty teams", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		// Create transaction mock
		dbMock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		// Prepare test data with empty teams
		gameDate := time.Now()
		gameReq := domain.GameStatsReq{
			Date:  gameDate,
			Teams: []domain.Team{},
		}

		// Test
		gameID, err := repo.Save(tx, gameReq)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, gameID)
	})
}

func TestRepository_Find(t *testing.T) {
	t.Run("find game stats", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		gameID := uuid.New().String()
		gameDate := time.Now().Format(time.DateTime)

		// Create mock data
		rows := sqlmock.NewRows([]string{
			"game_id", "player_id", "name", "date", "points", "rebounds", "assists",
			"steals", "blocks", "fouls", "turnovers", "minutes_played",
		}).
			AddRow(gameID, "player1", "LeBron James", gameDate, 24, 10, 8, 2, 1, 2, 3, 36).
			AddRow(gameID, "player2", "Anthony Davis", gameDate, 28, 12, 3, 1, 3, 2, 1, 34).
			AddRow(gameID, "player3", "Russell Westbrook", gameDate, 18, 7, 10, 3, 0, 3, 4, 32)

		// Set up expectations for the SELECT query
		dbMock.ExpectQuery("select g.game_id, g.player_id, p.name, g.date, g.points, g.rebounds, g.assists, g.steals, g.blocks, g.fouls, g.turnovers, g.minutes_played from game_stats g join players p").
			WithArgs(gameID).
			WillReturnRows(rows)

		// Test
		results, err := repo.Find(gameID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 3)

		// Check the first player's stats
		assert.Equal(t, "LeBron James", results[0].Name)
		assert.Equal(t, "player1", results[0].PlayerID)
		assert.Equal(t, 24, results[0].Points)
		assert.Equal(t, 10, results[0].Rebounds)
		assert.Equal(t, 8, results[0].Assists)
		assert.Equal(t, 2, results[0].Steals)
		assert.Equal(t, 1, results[0].Blocks)
		assert.Equal(t, 2, results[0].Fouls)
		assert.Equal(t, 3, results[0].Turnovers)
		assert.Equal(t, float64(36), results[0].MinutesPlayed)
	})

	t.Run("no game stats found", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		gameID := uuid.New().String()

		// Create empty result set
		rows := sqlmock.NewRows([]string{
			"game_id", "player_id", "name", "date", "points", "rebounds", "assists",
			"steals", "blocks", "fouls", "turnovers", "minutes_played",
		})

		// Set up expectations for the SELECT query
		dbMock.ExpectQuery("select g.game_id, g.player_id, p.name, g.date, g.points, g.rebounds, g.assists, g.steals, g.blocks, g.fouls, g.turnovers, g.minutes_played from game_stats g join players p").
			WithArgs(gameID).
			WillReturnRows(rows)

		// Test
		results, err := repo.Find(gameID)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("query error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		gameID := uuid.New().String()

		// Set up expectations for the SELECT query to fail
		dbMock.ExpectQuery("select g.game_id, g.player_id, p.name, g.date, g.points, g.rebounds, g.assists, g.steals, g.blocks, g.fouls, g.turnovers, g.minutes_played from game_stats g join players p").
			WithArgs(gameID).
			WillReturnError(sql.ErrConnDone)

		// Test
		results, err := repo.Find(gameID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("scan error", func(t *testing.T) {
		// Setup
		db, dbMock, _ := createMockDB(t)
		logger := zaptest.NewLogger(t)
		repo := NewRepo(db, logger)

		gameID := uuid.New().String()
		gameDate := "invalid-date-format"

		// Create mock data with invalid date format
		rows := sqlmock.NewRows([]string{
			"game_id", "player_id", "name", "date", "points", "rebounds", "assists",
			"steals", "blocks", "fouls", "turnovers", "minutes_played",
		}).
			AddRow(gameID, "player1", "LeBron James", gameDate, 24, 10, 8, 2, 1, 2, 3, 36)

		// Set up expectations for the SELECT query
		dbMock.ExpectQuery("select g.game_id, g.player_id, p.name, g.date, g.points, g.rebounds, g.assists, g.steals, g.blocks, g.fouls, g.turnovers, g.minutes_played from game_stats g join players p").
			WithArgs(gameID).
			WillReturnRows(rows)

		// Test
		results, err := repo.Find(gameID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, results)
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
