package forms

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Message struct {
	Id      int
	Name    string
	Email   string
	Subject string
	Message string
}

type MessageModel struct {
	DB *sql.DB
}

type User struct {
	Username string
	Password string
}

type UserModel struct {
	DB *sql.DB
}

func (m *MessageModel) Insert(name, email, subject, message string) error {
	query := "INSERT INTO form (name, email, subject, message) VALUES (?, ?, ?, ?)"
	_, err := m.DB.Exec(query, name, email, subject, message)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) CheckLogin(user, pass string) error {
	id := 1
	query := "SELECT * FROM users WHERE username = ? AND passwd = ?"
	if err := u.DB.QueryRow(query, user, pass).Scan(&id, &user, &pass); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	return nil
}

func (m *MessageModel) ShowMSGs() ([]*Message, error) {
	query := "SELECT * FROM form"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		msg := &Message{}
		id := 1
		if err := rows.Scan(&msg.Id, &msg.Email, &msg.Name, &msg.Subject, &msg.Message); err != nil {
			return messages, err
		}
		messages = append(messages, msg)
		id++
	}
	if err = rows.Err(); err != nil {
		return messages, err
	}
	return messages, nil
}
