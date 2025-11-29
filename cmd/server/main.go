package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sliitmozilla/auth/app/router"
	"github.com/sliitmozilla/auth/config"
)

func main() {
	c := config.GetConfig()
	r := chi.NewRouter()
	r.Mount("/", router.SetupRoutes())
	http.ListenAndServe(c.Host+":"+c.Port, r)
}
