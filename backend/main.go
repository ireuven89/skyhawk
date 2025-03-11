package main

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"skyhawk/backend/game/db"
	handler2 "skyhawk/backend/game/handler"
	"skyhawk/backend/game/usecase"
	goose "skyhawk/backend/goose"
	playerrepo "skyhawk/backend/player/db"
	"skyhawk/backend/redis"
	teamrepo "skyhawk/backend/team/db"
)

var migrationsDir = "/backend/goose/migrations"

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed initaiting service logger error %v", err)
	}

	//prepare DBs
	DB, err := goose.MustNewDB()
	if err != nil {
		log.Fatal(err)
	}

	//migrate database
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrationsDir, err := filepath.Abs(pwd + migrationsDir)
	if err != nil {
		panic(err)
	}
	migrationService := goose.New(DB, logger, migrationsDir)
	if err = migrationService.Run(); err != nil {
		panic(err)
	}

	redis, err := redis.MustNewRedis()

	if err != nil {
		log.Fatal(err)
	}
	//initiate service
	playerRepo := playerrepo.NewRepo(logger, DB, redis)
	teamRepo := teamrepo.New(DB, redis, logger)
	gameRepo := db.NewRepo(DB, logger)
	service := usecase.NewUseCase(logger, gameRepo, teamRepo, playerRepo)

	//handler
	handler := handler2.NewHandler(service, logger)

	e := echo.New()

	group := e.Group("api/v1")

	//game handler
	group.Add(http.MethodPost, "/games/log", handler.GameLogHandler)
	group.Add(http.MethodGet, "/games/:id", handler.GameStatsHandler)

	//player handler
	group.Add(http.MethodGet, "/players/season/:player_id", handler.PlayerSeasonStatsHandler)

	//team handler

	group.Add(http.MethodGet, "/teams/stats/season/:team_id", handler.TeamSeasonStatsHandler)

	log.Fatal(e.Start(":8080"))

}
