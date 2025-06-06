package dbx

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка ping БД: %v", err)
	}

	return nil
}
