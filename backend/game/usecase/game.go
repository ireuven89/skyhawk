package usecase

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"strings"
	"time"

	game_domain "skyhawk/backend/game/domain"
	player_domain "skyhawk/backend/player/domain"
	"skyhawk/backend/team/domain"
)

type PlayerRepository interface {
	SeasonStats(id string) (player_domain.PlayerSeasonStats, error)
	Save(ctx context.Context, tx *sql.Tx, player []player_domain.Player) (map[string]string, error)
}

type TeamRepository interface {
	Save(context context.Context, tx *sql.Tx, team domain.Team) (string, error)
	GetStats(id string) (domain.SeasonStats, error)
}

type GameRepository interface {
	Begin() (tx *sql.Tx, err error)
	Save(tx *sql.Tx, game game_domain.GameStatsReq) (string, error)
	Find(id string) ([]game_domain.GameStats, error)
}

type GameUseCase interface {
	GetGameStats(id string) ([]game_domain.GameStats, error)
	GetPlayerSeasonStats(id string) (player_domain.PlayerSeasonStats, error)
	GetTeamSeasonStats(id string) (domain.SeasonStats, error)
	LogGame(stats game_domain.GameStatsReq) (string, error)
}

type UseCase struct {
	gameRepo   GameRepository
	teamRepo   TeamRepository
	playerRepo PlayerRepository
	logger     *zap.Logger
}

const maxRetries = 3

func NewUseCase(logger *zap.Logger, gameRepo GameRepository, teamRepo TeamRepository, playerRepo PlayerRepository) *UseCase {

	return &UseCase{
		gameRepo:   gameRepo,
		teamRepo:   teamRepo,
		playerRepo: playerRepo,
		logger:     logger,
	}
}

func (s *UseCase) LogGame(stats game_domain.GameStatsReq) (string, error) {
	var id string

	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		id, err := s.attemptTransaction(stats)

		if err == nil {
			// Success
			return id, nil
		}

		// Check if it's a deadlock error
		if strings.Contains(err.Error(), "Deadlock found") {
			s.logger.Warn("Deadlock detected, retrying transaction",
				zap.Int("attempt", retryCount+1),
				zap.Int("maxRetries", maxRetries))

			// Add exponential backoff
			time.Sleep(time.Millisecond * 100 * time.Duration(retryCount+1))
			continue
		}

		// Not a deadlock error, return it
		return "", err
	}

	return id, nil
}

func (s *UseCase) attemptTransaction(stats game_domain.GameStatsReq) (string, error) {
	// Start transaction
	tx, err := s.gameRepo.Begin()
	if err != nil {
		s.logger.Error("UseCase.LogGame failed initiating transaction", zap.Error(err))
		return "", err
	}
	defer tx.Rollback()

	// Save all teams and update their IDs in stats
	for i := range stats.Teams {
		id, err := s.teamRepo.Save(context.Background(), tx, domain.Team{
			ID:   stats.Teams[i].ID,
			Name: stats.Teams[i].Name,
		})
		if err != nil {
			s.logger.Error("UseCase.LogGame failed logging teams", zap.Error(err))
			return "", err
		}
		stats.Teams[i].ID = id
	}

	// Prepare players
	var players []player_domain.Player
	for _, team := range stats.Teams {
		for i := range team.Players {
			players = append(players, player_domain.Player{
				Team: team.ID, // Now using the updated team ID
				Name: team.Players[i].Name,
			})
		}
	}

	// Insert players
	playerIdsMap, err := s.playerRepo.Save(context.Background(), tx, players)
	if err != nil {
		s.logger.Error("UseCase.LogGame failed processing players", zap.Error(err))
		return "", err
	}

	// Update player IDs in stats
	for i := range stats.Teams {
		for j := range stats.Teams[i].Players {
			playerName := stats.Teams[i].Players[j].Name
			if id, exists := playerIdsMap[playerName]; exists {
				stats.Teams[i].Players[j].ID = id
			}
		}
	}

	// insert game stats
	id, err := s.gameRepo.Save(tx, stats)

	if err != nil {
		s.logger.Error("failed saving game stats", zap.Error(err))
		return "", err

	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		s.logger.Error("failed committing changes", zap.Error(err))
		return "", err
	}

	return id, nil
}

func (s *UseCase) GetPlayerSeasonStats(id string) (player_domain.PlayerSeasonStats, error) {

	stats, err := s.playerRepo.SeasonStats(id)

	if err != nil {
		s.logger.Error("UseCase.GetPlayerSeasonStats failed fetching stats", zap.Error(err))
		return player_domain.PlayerSeasonStats{}, err
	}

	return stats, nil
}

func (s *UseCase) GetGameStats(id string) ([]game_domain.GameStats, error) {
	stats, err := s.gameRepo.Find(id)

	if err != nil {
		s.logger.Error("UseCase.GetGameStats failed fetching team stats", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

func (s *UseCase) GetTeamSeasonStats(id string) (domain.SeasonStats, error) {

	stats, err := s.teamRepo.GetStats(id)

	if err != nil {
		s.logger.Error("UseCase.GetTeamSeasonStats failed fetching team stats", zap.Error(err))
		return domain.SeasonStats{}, err
	}

	return stats, nil
}
