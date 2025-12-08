package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sliitmozilla/accounts/app/router"
)

// @title					 sliitmozilla Auth Service
// @description				 API documentation for the authentication service used across all sliitmozilla
// @version					 1.0
// @schemes					 https http
// @host					 accounts.sliitmozilla.org
// @basePath				 /api
// @contact.name			 Mozilla Campus Club of SLIIT
// @contact.url 			 https://www.sliitmozilla.org/contact/
// @contact.email 			 infosliitmcc@gmail.com
// @license.name 			 MPL-2.0
// @license.url 			 https://github.com/Mozilla-Campus-Club-of-SLIIT/accounts/blob/main/LICENSE
// @accept 					 json
// @produce 				 json
// @securityDefinitions.http AccessToken
// @scheme 					 bearer
// @bearerFormat 			 JWT
func Handler(w http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Mount("/api", router.SetupRoutes())
	r.ServeHTTP(w, req)
}
