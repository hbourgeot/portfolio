package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"henrry.online/internal/hangman"
	"henrry.online/internal/store"
)

func (folio *portfolio) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		folio.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamicMiddleware := alice.New(folio.sessionManager.LoadAndSave)

	// Portfolio handlers
	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(folio.home))
	router.Handler(http.MethodPost, "/send", dynamicMiddleware.ThenFunc(folio.SendMessage))
	router.Handler(http.MethodGet, "/login", dynamicMiddleware.ThenFunc(folio.Login))
	router.Handler(http.MethodPost, "/login", dynamicMiddleware.ThenFunc(folio.LoginPost))
	router.Handler(http.MethodGet, "/panel/:logged", dynamicMiddleware.ThenFunc(folio.ShowPanel))

	// API handlers
	router.Handler(http.MethodGet, "/api/get/hangman/:id", dynamicMiddleware.ThenFunc(hangman.GetHangman))
	router.Handler(http.MethodPost, "/api/new/hangman/:hint/:answer", dynamicMiddleware.ThenFunc(hangman.CreateHangman))

	// Store handlers
	router.Handler(http.MethodGet, "/store", dynamicMiddleware.ThenFunc(folio.StoreHome))
	router.Handler(http.MethodGet, "/store/login", dynamicMiddleware.ThenFunc(folio.StoreLogin))
	router.Handler(http.MethodPost, "/store/login", dynamicMiddleware.ThenFunc(store.LoginPost))
	router.Handler(http.MethodGet, "/store/products", dynamicMiddleware.ThenFunc(store.ShowProducts))
	router.Handler(http.MethodPost, "/store/products", dynamicMiddleware.ThenFunc(store.NewOrder))
	router.Handler(http.MethodGet, "/store/get-products", dynamicMiddleware.ThenFunc(store.GetProductforCart))
	router.Handler(http.MethodGet, "/store/admin/:logged", dynamicMiddleware.ThenFunc(store.AdminCRUD))
	router.Handler(http.MethodPost, "/store/admin/:logged/new", dynamicMiddleware.ThenFunc(store.AdminCreate))
	router.Handler(http.MethodPost, "/store/admin/:logged/mod", dynamicMiddleware.ThenFunc(store.AdminUpdate))
	router.Handler(http.MethodPost, "/store/admin/:logged/del", dynamicMiddleware.ThenFunc(store.AdminDelete))

	// Standard middlewares
	standardMiddleware := alice.New(folio.recoverPanic, folio.logRequest, folio.secureHeaders)

	return standardMiddleware.Then(router)
}
