package crud

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func makeCN() (*sql.DB, error) {
	conn := "postgres://henrrybrgt:Reyshell@full_db_postgres:5432/store"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
