package domain

type Player struct {
	ID   string
	Name string
	Team string
}

type PlayerSeasonStats struct {
	PlayerID         string
	PlayerName       string
	TeamID           string
	TeamName         string
	GamesPlayed      int
	AvgPoints        float64
	AvgRebounds      float64
	AvgAssists       float64
	AvgSteals        float64
	AvgBlocks        float64
	AvgFouls         float64
	AvgTurnovers     float64
	AvgMinutesPlayed float64
}
