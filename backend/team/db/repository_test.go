package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"skyhawk/backend/team/domain"
)

func TestUpsert(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Create an instance of Repo

	repo := &Repo{db: sqlx.NewDb(mockDB, "mysql"), logger: zaptest.NewLogger(t)}

	// Generate a random UUID for the team ID
	teamID := uuid.New().String()

	// Expected SQL query for INSERT
	mock.ExpectExec("INSERT INTO teams (id, name) VALUES (?,?) ON DUPLICATE KEY UPDATE name = ?").
		WithArgs(teamID, "Team A", "Team A").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate a successful insert

	// Call the Upsert function
	team := domain.Team{ID: "", Name: "Team A"}
	id, err := repo.Upsert(team)

	// Validate the result
	assert.NoError(t, err)
	assert.NotEmpty(t, id) // Assert that the id is generated

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamStats(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Create an instance of Repo
	repo := &Repo{db: sqlx.NewDb(mockDB, "mysql")}

	// Prepare mock data
	mockRow := sqlmock.NewRows([]string{"team_id", "team_name", "avg_rebounds", "avg_assists", "avg_steals", "avg_blocks", "avg_fouls", "avg_turnovers", "avg_minutes_played"}).
		AddRow("team-1", "Team A", 10.5, 5.5, 3.2, 2.1, 1.0, 3.0, 28.4)

	// Expected SQL query for SELECT
	mock.ExpectQuery("SELECT team_id, team_name, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played FROM player_season_stats WHERE player_id = ?").
		WithArgs("player-1").
		WillReturnRows(mockRow) // Simulate a row returned

	// Call the TeamStats function
	stats, err := repo.TeamStats("player-1")

	// Validate the result
	assert.NoError(t, err)
	assert.Equal(t, "team-1", stats.TeamID)
	assert.Equal(t, "Team A", stats.TeamName)
	assert.Equal(t, 10.5, stats.AvgRebounds)
	assert.Equal(t, 5.5, stats.AvgAssists)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
