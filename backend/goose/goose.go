package goose

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type Service interface {
	migrateDB() error
	Run() error
}

type MigrationService struct {
	db            *sqlx.DB
	logger        *zap.Logger
	migrationsDir string
}

func New(db *sqlx.DB, logger *zap.Logger, migrationsDir string) Service {

	return &MigrationService{
		db:            db,
		logger:        logger,
		migrationsDir: migrationsDir,
	}
}

// Run - this function migrates DB
func (ms *MigrationService) Run() error {

	if err := ms.migrateDB(); err != nil {
		ms.logger.Error("failed migrating DB", zap.Error(err))
		return err
	}

	return nil
}

// migrateDB - this function migrates the DB
func (ms *MigrationService) migrateDB() error {
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	start := time.Now()
	ms.logger.Info("starting migration...")
	if err := goose.Up(ms.db.DB, ms.migrationsDir); err != nil {
		ms.logger.Error("failed migration")
		return err
	}
	end := time.Now()
	ms.logger.Info("finished migration time to run: ", zap.String("dir", ms.migrationsDir), zap.Any("seconds", end.Second()-start.Second()))

	return nil
}
