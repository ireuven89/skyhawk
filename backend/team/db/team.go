package db

type Team struct {
	ID   string `goose:"id"`
	Name string `goose:"name"`
}

type SeasonStats struct {
	TeamID           string  `json:"team_id" goose:"team_id"`
	TeamName         string  `json:"team_name" goose:"team_name"`
	GamesPlayed      int     `json:"games_played" goose:"games_played"`
	AvgPoints        float64 `json:"avg_points" goose:"avg_points"`
	AvgRebounds      float64 `json:"avg_rebounds" goose:"avg_rebounds"`
	AvgAssists       float64 `json:"avg_assists" goose:"avg_assists"`
	AvgSteals        float64 `json:"avg_steals" goose:"avg_steals"`
	AvgBlocks        float64 `json:"avg_blocks" goose:"avg_blocks"`
	AvgFouls         float64 `json:"avg_fouls" goose:"avg_fouls"`
	AvgTurnovers     float64 `json:"avg_turnovers" goose:"avg_turnovers"`
	AvgMinutesPlayed float64 `json:"avg_minutes_played" goose:"avg_minutes_played"`
}
