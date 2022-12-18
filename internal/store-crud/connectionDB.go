package crud

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func makeCN() (*sql.DB, error) {
	conn := os.Getenv("TODOTECH_CONN_STRING")
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
