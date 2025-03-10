package domain

type Team struct {
	ID   string
	Name string
}

type SeasonStats struct {
	TeamID           string
	TeamName         string
	GamesPlayed      int     `json:"games_played"`
	AvgPoints        float64 `json:"avg_points" `
	AvgRebounds      float64 `json:"avg_rebounds" `
	AvgAssists       float64 `json:"avg_assists" `
	AvgSteals        float64 `json:"avg_steals" `
	AvgBlocks        float64 `json:"avg_blocks"`
	AvgFouls         float64 `json:"avg_fouls"`
	AvgTurnovers     float64 `json:"avg_turnovers"`
	AvgMinutesPlayed float64 `json:"avg_minutes_played"`
}
