package router

import (
	"github.com/go-chi/chi/v5"

	authHandlers "github.com/sliitmozilla/auth/app/handlers"
)

type AuthRoutes struct{}

func (b AuthRoutes) Routes() chi.Router {

	r := chi.NewRouter()

	r.Route("/", func(testRoutes chi.Router) {
		testRoutes.Post("/login", authHandlers.Login)
		testRoutes.Post("/logout", authHandlers.Logout)
	})

	return r
}
