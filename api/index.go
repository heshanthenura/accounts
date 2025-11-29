package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sliitmozilla/auth/app/router"
)

// @title SLIIT Mozilla Club Auth Service
// @description API documentation for the authentication service used across all sliitmozilla
// @version 1.0
// @schemes https http
// @host auth.sliitmozilla.org

// Entrypoint
func Handler(w http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Mount("/", router.SetupRoutes())
	r.ServeHTTP(w, req)
}
