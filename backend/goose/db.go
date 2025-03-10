package goose

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func MustNewDB() (*sqlx.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	url := fmt.Sprintf("root:%s@tcp(%s)/%s", password, host, "games_db")
	fmt.Println(url)
	db, err := sqlx.Connect("mysql", url)

	if err != nil {
		return nil, err
	}

	//ping
	if err = db.Ping(); err != nil {
		fmt.Printf(fmt.Sprintf("failed connecting to goose %v", err))
		return nil, err
	}

	return db, nil
}
