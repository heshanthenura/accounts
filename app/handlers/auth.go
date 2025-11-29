package handlers

import "net/http"

// @Tags        Auth
// @Summary     Login user
// @Description Endpoint to log in a user with credentials or session token
// @Accept      json
// @Produce     json
// @Param       username body string true "Username"
// @Param       password body string true "Password"
// @Success     200 {object} object "Login successful, returns session info"
// @Failure     400 {object} object "Invalid request or missing parameters"
// @Failure     401 {object} object "Invalid credentials"
// @Failure     500 {object} object "Server error"
// @Router      /login [POST]
func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login"))
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
