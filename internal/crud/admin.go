package crud

func ConfirmLogin(user, pass string) error {
	db, err := makeCN()
	if err != nil {
		return err
	}

	query := "SELECT username, password, country FROM customers WHERE username = $1 AND password = $2"
	row := db.QueryRow(query, user, pass)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}