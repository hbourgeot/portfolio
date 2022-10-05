package crud

import (
	_ "github.com/lib/pq"
)

func CreateClient(username, name, country string) error {
	db, err := makeCN()
	if err != nil {
		return err
	}

	query := "INSERT INTO customers (username,name,country) VALUES ($1,$2,$3)"
	result, err := db.Exec(query, username, name, country)
	if err != nil {

		return err
	}

	_, err = result.RowsAffected()
	if err != nil {

		return err
	}

	return nil
}
