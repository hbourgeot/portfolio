package hangman

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/hbourgeot/portfolio/internal/api"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetHangman(w http.ResponseWriter, r *http.Request) {
	db, err := openDB("root:Wini.h16b.@/hangman?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	quest := &api.QuestsModel{DB: db}

	enableCors(&w)
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "Method not allowed. %d", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/new/hangman" {
		return
	}

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		log.Fatal(err)
	}

	question := quest.GetQuestion(id)

	json.NewEncoder(w).Encode(question)
}

func CreateHangman(w http.ResponseWriter, r *http.Request) {
	db, err := openDB("root:Wini.h16b.@/hangman?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	quest := &api.QuestsModel{DB: db}
	enableCors(&w)

	params := httprouter.ParamsFromContext(r.Context())
	hint, answer := params.ByName("hint"), params.ByName("answer")
	fmt.Println(hint, answer)
	err = quest.Insert(hint, answer)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "Added hint: %s with answer: %s to the database!", hint, answer)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
