package model

import "time"

type Team struct {
	ID      string        `json:"-"`
	Name    string        `json:"name"`
	Players []PlayerStats `json:"players¬"`
}

type GameStats struct {
	ID    string    `json:"ID"`
	Date  time.Time `json:"date"`
	Name  string    `json:"name"`
	Teams []Team    `json:"teams"`
}

type PlayerStats struct {
	ID            string  `json:"ID"`
	PLayer        string  `json:"PLayer"`
	PlayerName    string  `json:"playerName"`
	Points        int     `json:"points"`
	Rebounds      int     `json:"rebounds"`
	Assists       int     `json:"assists"`
	Steals        int     `json:"steals"`
	Blocks        int     `json:"blocks"`
	Fouls         int     `json:"fouls"`
	Turnovers     int     `json:"turnovers"`
	MinutesPlayed float64 `json:"minutes_played"`
}
