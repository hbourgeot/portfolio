package main

import (
	"log"
	"net/http"

	"github.com/hbourgeot/portfolio/internal/api"
	forms "github.com/hbourgeot/portfolio/internal/models"
	"github.com/hbourgeot/portfolio/internal/store"
)

func main() {
	forms.ConnectToDB()
	fileServer := http.FileServer(http.Dir("go/src/portfolio/ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Portfolio handlers
	http.HandleFunc("/", home)
	http.HandleFunc("/not-found", notFound)
	http.HandleFunc("/admin", LoginPage)
	http.HandleFunc("/admin/panel", ShowPanel)
	http.HandleFunc("/submit", SendForm)
	http.HandleFunc("/login", Login)
	log.Printf("Starting server on port :4040")

	// API handlers
	fileServer = http.FileServer(http.Dir("go/src/portfolio/internal/api/static"))
	http.Handle("/api", fileServer)
	http.HandleFunc("/api/get/hangman", api.GetHangman)
	http.HandleFunc("/api/new/hangman", api.CreateHangman)

	// Store handlers
	fileServer = http.FileServer(http.Dir("go/src/portfolio/internal/store/ui/static"))
	http.Handle("/static/store", fileServer)
	http.HandleFunc("/store", store.Home)
	http.HandleFunc("/store/products", store.ShowProducts)
	http.HandleFunc("/store/products/cart", store.GetProductforCart)
	http.HandleFunc("/store/neworder", store.NewOrder)
	http.HandleFunc("/store/admin", store.Login)
	http.HandleFunc("/store/admin/panel", store.AdminCRUD)
	http.HandleFunc("/store/admin/panel/ins", store.AdminCreate)
	http.HandleFunc("/store/admin/panel/upda", store.AdminUpdate)
	http.HandleFunc("/store/admin/panel/del", store.AdminDelete)
	http.HandleFunc("/store/submit", store.LoginAdm)

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
