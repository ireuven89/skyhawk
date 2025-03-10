package db

import "time"

type GameStatsDB struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	PlayerID      string    `db:"player_id"`
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
