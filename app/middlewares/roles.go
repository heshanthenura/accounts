package middlewares

import (
	"net/http"
	"slices"

	"github.com/sliitmozilla/accounts/db/models"
	"github.com/sliitmozilla/accounts/helpers"
)

func RequireRoles(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(UserContext{}) == nil {
				helpers.Response(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			ctxUser, ok := r.Context().Value(UserContext{}).(*models.UserModel)
			if ctxUser == nil || !ok {
				helpers.Response(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			userRoles := ctxUser.Roles
			for _, role := range requiredRoles {
				if !slices.Contains(userRoles, role) {
					helpers.Response(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
