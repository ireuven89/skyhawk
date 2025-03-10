package db

import "time"

type PlayerStatsDB struct {
	PlayerID      string    `db:"player_id"`
	PlayerName    string    `db:"player_name"`
	GameID        string    `db:"game_id"`
	Date          time.Time `db:"date"`
	Points        int       `db:"points"`
	Rebounds      int       `db:"rebounds"`
	Assists       int       `db:"assists"`
	Steals        int       `db:"steals"`
	Blocks        int       `db:"blocks"`
	Fouls         int       `db:"fouls"`
	Turnovers     int       `db:"turnovers"`
	MinutesPlayed float64   `db:"minutes_played"`
}

type Player struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Team string `db:"team_id"`
}

type PlayerSeasonStats struct {
	PlayerID         string  `db:"player_id"`
	PlayerName       string  `db:"player_name"`
	TeamID           string  `db:"team_id"`
	TeamName         string  `db:"team_name"`
	GamesPlayed      int     `db:"games_played"`
	AvgPoints        float64 `db:"avg_points"`
	AvgRebounds      float64 `db:"avg_rebounds"`
	AvgAssists       float64 `db:"avg_assists"`
	AvgSteals        float64 `db:"avg_steals"`
	AvgBlocks        float64 `db:"avg_blocks"`
	AvgFouls         float64 `db:"avg_fouls"`
	AvgTurnovers     float64 `db:"avg_turnovers"`
	AvgMinutesPlayed float64 `db:"avg_minutes_played"`
}
