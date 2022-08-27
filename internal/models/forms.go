package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Message struct {
	Name    string
	Email   string
	Title   string
	Message string
}

var db *sql.DB

func makeConnection() {
	connStr := "host=localhost port=5432 user=postgres password=Reyshell dbname=todophone sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := db.Ping(); err != nil {
		log.Fatalln("Error al conectar a DB " + err.Error())
	}
}

func ConnectToDB() {
	makeConnection()
	fmt.Println("Conectado")
	defer db.Close()
}

func Insert(name, email, subject, message string) error {
	makeConnection()
	query := "INSERT INTO formulario (nombre, correo, subject, mensaje) VALUES ($1, $2, $3, $4)"
	_, e := db.Exec(query, name, email, subject, message)
	if e != nil {
		defer db.Close()
		return e
	}
	defer db.Close()
	return nil
}

func CheckLogin(user, pass string) (bool, error) {
	makeConnection()
	id := 1
	query := "SELECT * FROM users WHERE username = $1 AND passwd = $2"
	if err := db.QueryRow(query, user, pass).Scan(&id, &user, &pass); err != nil {
		if err == sql.ErrNoRows {
			defer db.Close()
			return false, err
		}
		defer db.Close()
		return false, err
	}
	defer db.Close()
	return true, nil
}

func ShowMSGs() ([]Message, error) {
	makeConnection()
	query := "SELECT * FROM formulario"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message
		id := 1
		if err := rows.Scan(&id, &msg.Name, &msg.Email, &msg.Title, &msg.Message); err != nil {
			defer db.Close()
			return messages, err
		}
		messages = append(messages, msg)
		id++
	}
	if err = rows.Err(); err != nil {
		defer db.Close()
		return messages, err
	}
	defer db.Close()
	return messages, nil
}
