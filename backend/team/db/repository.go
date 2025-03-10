package db

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"skyhawk/backend/team/domain"
)

type Repository interface {
	Save(tx *sql.Tx, team domain.Team) (string, error)
	Find(id string) (domain.Team, error)
	GetStats(id string) (domain.SeasonStats, error)
}

type Repo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func New(db *sqlx.DB, logger *zap.Logger) Repository {

	return &Repo{
		db:     db,
		logger: logger,
	}
}

func (r *Repo) Save(tx *sql.Tx, team domain.Team) (string, error) {
	if team.ID == "" {
		team.ID = uuid.New().String()
	}

	row, err := tx.Query("INSERT INTO teams (id, name) VALUES (?,?) on duplicate key update name=values(name)", team.ID, team.Name)
	if err != nil || row.Err() != nil {
		r.logger.Error("failed inserting team", zap.Error(err))
		return "", err
	}

	return team.ID, nil
}

func (r *Repo) Find(id string) (domain.Team, error) {
	var teamDB Team

	row := r.db.QueryRow("select id, name  from teams where name = ?", id)

	if err := row.Scan(teamDB); err != nil {
		return domain.Team{}, err
	}

	return domain.Team{
		ID:   teamDB.ID,
		Name: teamDB.Name,
	}, nil
}

func (r *Repo) GetStats(id string) (domain.SeasonStats, error) {
	var seasonDB SeasonStats

	row := r.db.QueryRow("select team_id, team_name, avg_rebounds, avg_assists, avg_steals, avg_blocks, avg_fouls, avg_turnovers, avg_minutes_played  from player_season_stats where player_id = ?", id)

	if err := row.Scan(seasonDB); err != nil {
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
