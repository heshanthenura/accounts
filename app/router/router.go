package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/sliitmozilla/auth/docs"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	dir, _ := os.Getwd()

	r.Use(middleware.Logger)

	r.Mount("/", AuthRoutes{}.Routes())

	// Swagger docs
	fs := http.FileServer(http.Dir(filepath.Join(dir, "docs")))
	r.Handle("/docs/swagger.json", http.StripPrefix("/docs", fs))
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	return r
}
