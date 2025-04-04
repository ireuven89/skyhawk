package domain

import "time"

type GameStatsReq struct {
	ID     string    `json:"id"`
	GameID string    `json:"gameID"`
	Date   time.Time `json:"date"`
	Teams  []Team    `json:"teams"`
}

type Team struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Players []Player `json:"players"`
}

type Player struct {
	ID            string  `json:"id"`
	Name          string  `json:"player_name"`
	Points        int     `json:"points"`
	Rebounds      int     `json:"rebounds"`
	Assists       int     `json:"assists"`
	Steals        int     `json:"steals"`
	Blocks        int     `json:"blocks"`
	Fouls         int     `json:"fouls"`
	Turnovers     int     `json:"turnovers"`
	MinutesPlayed float64 `json:"minutes_played"`
}

type GameStats struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`
	Name          string    `json:"name"`
	PlayerID      string    `json:"player_id"`
	Points        int       `json:"points"`
	Rebounds      int       `json:"rebounds"`
	Assists       int       `json:"assists"`
	Steals        int       `json:"steals"`
	Blocks        int       `json:"blocks"`
	Fouls         int       `json:"fouls"`
	Turnovers     int       `json:"turnovers"`
	MinutesPlayed float64   `json:"minutes_played"`
}
