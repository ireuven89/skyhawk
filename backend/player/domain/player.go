package domain

type Player struct {
	ID   string
	Name string
	Team string
}

type PlayerSeasonStats struct {
	PlayerID         string  `json:"player_id"`
	PlayerName       string  `json:"player_name"`
	TeamID           string  `json:"team_id"`
	TeamName         string  `json:"team_name"`
	GamesPlayed      int     `json:"games_played"`
	AvgPoints        float64 `json:"avg_points"`
	AvgRebounds      float64 `json:"avg_rebounds"`
	AvgAssists       float64 `json:"avg_assists"`
	AvgSteals        float64 `json:"avg_steals"`
	AvgBlocks        float64 `json:"avg_blocks"`
	AvgFouls         float64 `json:"avg_fouls"`
	AvgTurnovers     float64 `json:"avg_turnovers"`
	AvgMinutesPlayed float64 `json:"avg_minutes_played"`
}
