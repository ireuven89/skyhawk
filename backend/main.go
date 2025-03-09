package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"skyhawk/backend/game/db"
	handler2 "skyhawk/backend/game/handler"
	"skyhawk/backend/game/usecase"
	goose "skyhawk/backend/goose"
	playerrepo "skyhawk/backend/player/db"
	teamrepo "skyhawk/backend/team/db"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed initaiting service logger error %v", err)
	}
	DB, err := goose.MustNewDB()
	if err != nil {
		log.Fatal(err)
	}

	//migrate goose
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrationsDir, err := filepath.Abs(pwd + "/backend/goose/migrations")
	if err != nil {
		panic(err)
	}
	migrationService := goose.New(DB, logger, migrationsDir)
	if err = migrationService.Run(); err != nil {
		panic(err)
	}

	playerRepo := playerrepo.NewRepo(logger, DB)
	teamRepo := teamrepo.New(DB, logger)
	gameRepo := db.NewRepo(DB, logger)
	service := usecase.NewUseCase(logger, gameRepo, teamRepo, playerRepo)

	handler := handler2.NewHandler(service, logger)

	e := echo.New()

	group := e.Group("api/v1")

	//game handler
	//	gamesGroup := api.Group("games")
	group.Add(http.MethodPost, "/games/log", handler.GameStatsHandler)
	group.Add(http.MethodGet, "/games/:id", handler.GameStatsHandler)

	//player handler
	//	playerGroup := api.Group("players")
	group.Add(http.MethodGet, "/players//:id/game/:id", handler.PlayerStatsHandler)
	group.Add(http.MethodGet, "/players/season/:player_id", handler.PlayerSeasonStatsHandler)

	//team handler
	//	teamGroup := api.Group("teams")
	group.Add(http.MethodGet, "/teams/stats/:team_id/:game_id", handler.TeamStatsStatsHandler)
	group.Add(http.MethodGet, "/teams/stats/season/:team_id", handler.TeamSeasonStatsHandler)

	log.Fatal(e.Start(":8080"))

}
