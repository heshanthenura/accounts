package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sliitmozilla/accounts/app/middlewares"
	"github.com/sliitmozilla/accounts/db/models"
	apiErrors "github.com/sliitmozilla/accounts/errors"
	"github.com/sliitmozilla/accounts/helpers"
)

// @tags        Auth
// @summary     Get current session
// @description Get the current session using the access token
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /session [GET]
func GetSession(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(middlewares.UserContext{}).(*models.UserModel)
	helpers.Response(w, http.StatusOK, map[string]any{
		"id":    u.ID.String(),
		"roles": u.Roles,
	})
}

// @tags        Auth
// @summary		Initiate the authentication flow
// @description Initiate the authentication flow with the auth service. \
// @description	Any external service should visit this route with a valid redirect \
// @description	If the user is already logged in with the auth service, the auth \
// @description	service will redirect the user back to the provided url with a \
// @description	temporary token - that should be used to complete the authentication \
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /authorize [GET]
func Authorize(w http.ResponseWriter, r *http.Request) {}

// @tags        Auth
// @summary     Login user
// @description Endpoint to log in a user with credentials or session token
// @accept      json
// @produce     json
// @param       username body string true "Username"
// @param       password body string true "Password"
// @success     200 {object} object "Login successful, returns session info"
// @failure     400 {object} object "Invalid request or missing parameters"
// @failure     401 {object} object "Invalid credentials"
// @failure     500 {object} object "Server error"
// @router      /login [POST]
func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		helpers.Response(w, http.StatusBadRequest, "Invalid or empty body")
		return
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		helpers.Response(w, http.StatusBadRequest, "email and password cannot be empty")
	}

	accessToken, refreshToken, err := models.UserModel{}.Login(requestBody.Email, requestBody.Password)
	if err != nil {
		if _, ok := err.(apiErrors.NotFoundError); ok {
			helpers.Response(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		log.Println(err.Error())
		helpers.Response(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "refreshToken",
		Value:   refreshToken,
		Expires: time.Now().Add(30 * 24 * time.Hour),
	})
	helpers.Response(w, http.StatusOK, map[string]string{"token": accessToken})
}

// @Tags        Auth
// @Summary     Logout user
// @Description Endpoint to log out a user and invalidate their session
// @Accept      json
// @Produce     json
// @Success     200 {object} object "Logout successful"
// @Failure     401 {object} object "Unauthorized / session not found"
// @Failure     500 {object} object "Server error"
// @Router      /logout [POST]
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}

func getAccessTokenFromRefresh(refreshToken string) (string, error) {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString := refreshToken
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	}); err != nil {
		return "", err
	}
	id := uuid.FromStringOrNil(claims["id"].(string))
	if id == uuid.Nil {
		return "", errors.New("invalid token")
	}

	u, err := models.UserModel{}.GetUserByID(id)
	if err != nil {
		return "", err
	}

	accessToken, _, err := helpers.GenerateTokens(u.ID.String(), u.Name, u.Email, u.Roles)
	return accessToken, err
}

// @Tags        Auth
// @Summary     Refresh acess token
// @Description Refresh the access token with the refresh token
// @Accept      json
// @Produce     json
// @Failure     500 {object} object "Server error"
// @Router      /token/refresh [POST]
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// todo: right now you can probably pass an access token and this endpoint would still work
	token, err := r.Cookie("refreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.Response(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		log.Println(err.Error())
		helpers.Response(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	accessToken, err := getAccessTokenFromRefresh(token.Value)
	if err != nil {
		if err.Error() == "invalid token" {
			helpers.Response(w, http.StatusUnauthorized, err.Error())
			return
		} else if _, ok := err.(apiErrors.NotFoundError); ok {
			helpers.Response(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		helpers.Response(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	helpers.Response(w, http.StatusOK, map[string]string{"token": accessToken})
}
