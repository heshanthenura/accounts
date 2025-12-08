package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/gofrs/uuid"
	middlewares "github.com/sliitmozilla/accounts/app/middlewares"
	models "github.com/sliitmozilla/accounts/db/models"
	errors "github.com/sliitmozilla/accounts/errors"
	helpers "github.com/sliitmozilla/accounts/helpers"
)

// @tags        Users
// @summary		List all users
// @description List all users. Admin only route
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users [GET]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.UserModel{}.SelectAll()
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.Response(w, http.StatusOK, users)
}

// @tags        Users
// @summary		Create user
// @description Create user
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var u models.UserModel
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helpers.Response(w, http.StatusBadRequest, "Invalid or empty body")
		return
	}
	if _, err := u.Insert(); err != nil {
		if ve, ok := err.(errors.ValidationError); ok {
			helpers.Response(w, http.StatusBadRequest, ve.Error())
			return
		} else if ve, ok := err.(errors.DuplicateError); ok {
			helpers.Response(w, http.StatusConflict, ve.Error())
			return
		}
		log.Println(err.Error())
		helpers.Response(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	helpers.Response(w, http.StatusCreated, http.StatusText(http.StatusCreated))
}

// @tags        Users
// @summary		Get current user
// @description Get current user
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/me [GET]
func GetMe(w http.ResponseWriter, r *http.Request) {
	ctxUser := r.Context().Value(middlewares.UserContext{}).(*models.UserModel)
	u, err := models.UserModel{}.GetUserByID(ctxUser.ID)
	if err != nil {
		if ve, ok := err.(errors.NotFoundError); ok {
			helpers.Response(w, http.StatusBadRequest, ve.Error())
			return
		}
		log.Println(err)
		helpers.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	helpers.Response(w, http.StatusOK, u)
}

// @tags        Users
// @summary		Update current user
// @description Update current user details such as username, social links, etc.
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/me [PATCH]
func UpdateMe(w http.ResponseWriter, r *http.Request) {}

// @tags        Users
// @summary		Change password of current user
// @description Change password of current user
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/me/password [PATCH]
func ChangePassword(w http.ResponseWriter, r *http.Request) {}

// @tags        Users
// @summary		Get a specific user
// @description Get a specific user
// @description If the user has set the profile visibility as public then it is \
// @description free to view by anyone. If the profile is set to be private, then \
// @description regular users will not gain the access to this resource.
// @description However, admin users can still bypass this protection
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/{id} [GET]
func GetUser(w http.ResponseWriter, r *http.Request) {
	requestedUserId := r.PathValue("id")
	requestedUserUuid := uuid.FromStringOrNil(requestedUserId)
	// check if the request context has enough permissions (self view or admin view)
	ctxUser := r.Context().Value(middlewares.UserContext{}).(*models.UserModel)
	if !slices.Contains(ctxUser.Roles, "admin") && requestedUserUuid != ctxUser.ID {
		helpers.Response(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	if requestedUserUuid == uuid.Nil {
		helpers.Response(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	u, err := models.UserModel{}.GetUserByID(requestedUserUuid)
	if err != nil {
		if ve, ok := err.(errors.NotFoundError); ok {
			helpers.Response(w, http.StatusNotFound, ve.Error())
			return
		}
		log.Println(err)
		helpers.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	helpers.Response(w, http.StatusOK, u)
}

// @tags        Users
// @summary		Update a user
// @description Update a user. Admin only route
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/{id} [PATCH]
func UpdateUser(w http.ResponseWriter, r *http.Request) {}

// @tags        Users
// @summary		Add a role to the user
// @description Add a role to the user
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/{id}/roles [POST]
func AddRole(w http.ResponseWriter, r *http.Request) {
	requestedUserId := r.PathValue("id")
	requestedUserUuid := uuid.FromStringOrNil(requestedUserId)
	role := models.RoleModel{}
	if requestedUserUuid == uuid.Nil {
		helpers.Response(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		helpers.Response(w, http.StatusBadRequest, "Invalid or empty body")
		return
	}
	if role.Name == "" {
		helpers.Response(w, http.StatusBadRequest, "role cannot be empty")
		return
	}
	u := models.UserModel{ID: requestedUserUuid}
	_, err := u.InsertRole(role.Name)
	if err != nil {
		if ve, ok := err.(errors.NotFoundError); ok {
			helpers.Response(w, http.StatusNotFound, ve.Error())
			return
		} else if ve, ok := err.(errors.DuplicateError); ok {
			helpers.Response(w, http.StatusConflict, ve.Error())
			return
		}
		log.Println(err.Error())
		helpers.Response(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	helpers.Response(w, http.StatusOK, http.StatusText(http.StatusOK))
}

// @tags        Users
// @summary		Remove an existing role from the user
// @description Remove an existing role from the user
// @security	AccessToken
// @accept      json
// @produce     json
// @failure     500 {object} object "Server error"
// @router      /users/{id}/roles/{role} [DELETE]
func RemoveRole(w http.ResponseWriter, r *http.Request) {
	requestedUserId := r.PathValue("id")
	requestedUserUuid := uuid.FromStringOrNil(requestedUserId)
	ctxUser := r.Context().Value(middlewares.UserContext{}).(*models.UserModel)
	role := r.PathValue("role")
	if requestedUserUuid == uuid.Nil {
		helpers.Response(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	if role == "admin" && requestedUserId == ctxUser.ID.String() {
		helpers.Response(w, http.StatusForbidden, "Cannot remove admin from yourself")
		return
	}
	u := models.UserModel{ID: requestedUserUuid}
	rows, err := u.RemoveRole(role)
	if rows == 0 {
		helpers.Response(w, http.StatusNotFound, "User or role not found")
		return
	}
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	helpers.Response(w, http.StatusOK, http.StatusText(http.StatusOK))
}
