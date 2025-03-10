package goose

import (
	"errors"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Mock DB connection
type MockDB struct {
	mock.Mock
}

// Mock goose.Up function
func (m *MockDB) MigrateUp(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

// MigrationService with dependency injection
type MockMigrationService struct {
	db            *MockDB
	migrationsDir string
	logger        *zap.Logger
}

// Migrate function using the mock DB
func (ms *MockMigrationService) migrateDB() error {
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	start := time.Now()
	ms.logger.Info("starting migration...")

	if err := ms.db.MigrateUp(ms.migrationsDir); err != nil {
		ms.logger.Error("failed migration")
		return err
	}

	end := time.Now()
	ms.logger.Info("finished migration",
		zap.String("dir", ms.migrationsDir),
		zap.Int("duration_seconds", end.Second()-start.Second()))

	return nil
}

// Test Migration Success
func TestMigrateDB_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockDB := new(MockDB)

	// Expect MigrateUp to be called and return no error
	mockDB.On("MigrateUp", "migrations").Return(nil)

	ms := &MockMigrationService{
		db:            mockDB,
		migrationsDir: "migrations",
		logger:        logger,
	}

	err := ms.migrateDB()
	require.NoError(t, err) // Ensure no error

	mockDB.AssertCalled(t, "MigrateUp", "migrations") // Ensure migration was attempted
}

// Test Migration Failure
func TestMigrateDB_Failure(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockDB := new(MockDB)

	// Expect MigrateUp to return an error
	mockDB.On("MigrateUp", "migrations").Return(errors.New("migration failed"))

	ms := &MockMigrationService{
		db:            mockDB,
		migrationsDir: "migrations",
		logger:        logger,
	}

	err := ms.migrateDB()
	require.Error(t, err) // Ensure error is returned

	mockDB.AssertCalled(t, "MigrateUp", "migrations") // Ensure migration was attempted
}
