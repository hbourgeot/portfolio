package main

import (
	"log"
	"net/http"

	forms "github.com/hbourgeot/portfolio/internal/models"
)

func main() {
	forms.ConnectToDB()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", home)
	http.HandleFunc("/not-found", notFound)
	http.HandleFunc("/admin", LoginPage)
	http.HandleFunc("/admin/panel", ShowPanel)
	http.HandleFunc("/submit", SendForm)
	http.HandleFunc("/login", Login)
	log.Printf("Starting server on port :4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
