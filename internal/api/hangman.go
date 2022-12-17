package api

import (
	"database/sql"
	"log"
)

type Quests struct {
	ID           int
	Hint         string
	Answer       string
	AnswerLenght int
}

type QuestsModel struct {
	DB *sql.DB
}

func (q *QuestsModel) GetQuestion(id int) []Quests {
	query := "SELECT hint, answer, answer_lenght FROM quests WHERE id = ?"
	row := q.DB.QueryRow(query, id)
	defer q.DB.Close()
	if row.Err() != nil {
		log.Fatal(row.Err())
		return nil
	}

	var quest Quests
	err := row.Scan(&quest.Hint, &quest.Answer, &quest.AnswerLenght)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var hangman []Quests
	hangman = append(hangman, quest)

	return hangman
}

func (q *QuestsModel) Insert(hint string, answer string) error {
	query := "INSERT INTO quests (hint, answer, answer_lenght) VALUES (?, ?, ?)"

	_, err := q.DB.Exec(query, hint, answer, len(answer))
	defer q.DB.Close()
	if err != nil {
		return err
	}

	return nil
}
