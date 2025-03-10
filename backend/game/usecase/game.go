package usecase

import (
	"database/sql"

	"go.uber.org/zap"

	game_domain "skyhawk/backend/game/domain"
	player_domain "skyhawk/backend/player/domain"
	"skyhawk/backend/team/domain"
)

type PlayerRepository interface {
	SeasonStats(id string) (player_domain.PlayerSeasonStats, error)
	Save(tx *sql.Tx, player player_domain.Player) (string, error)
}

type TeamRepository interface {
	Save(tx *sql.Tx, team domain.Team) (string, error)
	GetStats(id string) (domain.SeasonStats, error)
}

type GameRepository interface {
	Begin() (tx *sql.Tx, err error)
	Save(tx *sql.Tx, game game_domain.GameStats) (string, error)
	Find(id string) (game_domain.GameStats, error)
}

type GameUseCase interface {
	FindGame(id string) (game_domain.GameStats, error)
	FindPlayer(id string) (player_domain.PlayerSeasonStats, error)
	FindTeam(id string) (domain.SeasonStats, error)
	LogGame(stats game_domain.GameStats) (string, error)
}

type UseCase struct {
	gameRepo   GameRepository
	teamRepo   TeamRepository
	playerRepo PlayerRepository
	logger     *zap.Logger
}

func NewUseCase(logger *zap.Logger, gameRepo GameRepository, teamRepo TeamRepository, playerRepo PlayerRepository) *UseCase {

	return &UseCase{
		gameRepo:   gameRepo,
		teamRepo:   teamRepo,
		playerRepo: playerRepo,
		logger:     logger,
	}
}

func (s *UseCase) LogGame(stats game_domain.GameStats) (string, error) {
	tx, err := s.gameRepo.Begin()

	if err != nil {
		s.logger.Error("failed transaction")
		return "", err
	}

	for i := range stats.Teams {
		id, err := s.teamRepo.Save(tx, domain.Team{ID: stats.Teams[i].ID, Name: stats.Teams[i].Name})
		if err != nil {
			s.logger.Error("failed logging game", zap.Error(err))
			tx.Rollback()
			return "", err
		}
		stats.Teams[i].ID = id
	}

	for _, team := range stats.Teams {
		for i := range team.Players {
			id, err := s.playerRepo.Save(tx, player_domain.Player{ID: team.Players[i].ID, Name: team.Players[i].Name, Team: team.ID})
			if err != nil {
				s.logger.Error("failed inserting players", zap.Error(err))
				tx.Rollback()
				return "", err
			}
			team.Players[i].ID = id
		}
	}

	id, err := s.gameRepo.Save(tx, stats)

	if err != nil {
		s.logger.Error("failed logging game stats", zap.Error(err))
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return id, nil
}

func (s *UseCase) FindPlayer(id string) (player_domain.PlayerSeasonStats, error) {

	stats, err := s.playerRepo.SeasonStats(id)

	if err != nil {
		return player_domain.PlayerSeasonStats{}, err
	}

	return stats, nil
}

func (s *UseCase) FindGame(id string) (game_domain.GameStats, error) {
	stats, err := s.gameRepo.Find(id)

	if err != nil {
		return game_domain.GameStats{}, err
	}

	return stats, nil
}

func (s *UseCase) FindTeam(id string) (domain.SeasonStats, error) {

	stats, err := s.teamRepo.GetStats(id)

	if err != nil {

		return domain.SeasonStats{}, err
	}

	return stats, nil
}
