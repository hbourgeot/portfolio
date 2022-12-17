package crud

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func makeCN() (*sql.DB, error) {
	conn := "root:Wini.h16b.@/todotech?parseTime=true"
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
