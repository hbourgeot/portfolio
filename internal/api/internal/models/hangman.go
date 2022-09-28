package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Quests struct {
	ID           int
	Hint         string
	Answer       string
	AnswerLenght int
}

func makeCN() (*sql.DB, error) {
	connection := "postgres://henrrybrgt:Reyshell@full_db_postgres:5432/hangman"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	} else if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected")
	return db, nil
}

func GetQuestion(id int) []Quests {
	db, err := makeCN()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	query := "SELECT hint, answer, answer_lenght FROM quests WHERE id = $1"
	row := db.QueryRow(query, id)
	defer db.Close()
	if row.Err() != nil {
		log.Fatal(err)
		return nil
	}

	var quest Quests
	err = row.Scan(&quest.Hint, &quest.Answer, &quest.AnswerLenght)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var hangman []Quests
	hangman = append(hangman, quest)

	return hangman
}

func Insert(hint string, answer string) error {
	db, err := makeCN()
	if err != nil {
		log.Fatal(err)
		return err
	}

	query := "INSERT INTO quests (hint, answer, answer_lenght) VALUES ($1,$2,$3)"

	_, err = db.Exec(query, hint, answer, len(answer))
	defer db.Close()
	if err != nil {
		return err
	}

	return nil
}
