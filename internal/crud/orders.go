package crud

import (
	"database/sql"
	"fmt"
	"time"
)

func GenerateOrder(username string, cart []int) error {
	db, err := makeCN()

	t := time.Now()
	date := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	query := "INSERT INTO cart (username, product, order_gnrted) VALUES ($1,$2,$3)"
	var result sql.Result

	for i := range cart {
		result, err = db.Exec(query, username, cart[i], date)
	}

	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
