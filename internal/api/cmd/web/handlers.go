package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gamesapi/hangman/internal/models"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetHangman(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "Method not allowed. %d", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/new/hangman" {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}

	quest := models.GetQuestion(id)

	json.NewEncoder(w).Encode(quest)
}

func CreateHangman(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	hint, answer := r.URL.Query().Get("hint"), r.URL.Query().Get("answer")
	fmt.Println(hint, answer)
	err := models.Insert(hint, answer)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "Added hint: %s to the database!", hint)
}
