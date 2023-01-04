package main

import (
	"fmt"
	"net/http"
)

// Middleware to define secure headers
func (folio *portfolio) secureHeaders(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		nextHandler.ServeHTTP(w, r)
	})
}

func (folio *portfolio) logRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		folio.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		nextHandler.ServeHTTP(w, r)
	})
}

func (folio *portfolio) recoverPanic(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				folio.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		nextHandler.ServeHTTP(w, r)
	})
}
