package db

type Team struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type SeasonStats struct {
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
