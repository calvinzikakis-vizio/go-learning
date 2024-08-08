package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sqlite/authenticate"
	"sqlite/controllers"
	"sqlite/models"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// create a new user model from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Request body not valid JSON.", http.StatusBadRequest)
		return
	}

	// hash the password
	hash, err := authenticate.HashPassword([]byte(user.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a new user database model
	userDB := models.UserDatabase{
		Username: user.Username,
		Hash:     string(hash),
	}

	// insert the user into the database
	err = controllers.StorePasswordInfo(db, userDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a token for the user
	tokenString, err := authenticate.CreateToken([]byte(user.Username))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create response object
	createUserResponse := models.CreateUserResponse{
		Username: user.Username,
		Token:    tokenString,
	}

	// write the response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createUserResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// create a new user model from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Request body not valid JSON.", http.StatusBadRequest)
		return
	}

	// get the user from the database
	userDB, err := controllers.GetPasswordInfo(db, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// compare the password
	err = authenticate.CompareHashAndPasswords([]byte(userDB.Hash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
		return
	}

	// create a token for the user
	tokenString, err := authenticate.CreateToken([]byte(user.Username))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create response object
	createUserResponse := models.CreateUserResponse{
		Username: user.Username,
		Token:    tokenString,
	}

	// write the response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createUserResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// create a new user model from the request body
	var user models.UserChangePassword
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Request body not valid JSON.", http.StatusBadRequest)
		return
	}

	// validate the old password
	userDB, err := controllers.GetPasswordInfo(db, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = authenticate.CompareHashAndPasswords([]byte(userDB.Hash), []byte(user.OldPassword))
	if err != nil {
		http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
		return
	}

	// hash the new password
	hash, err := authenticate.HashPassword([]byte(user.NewPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update the password in the database
	userDB.Hash = string(hash)
	err = controllers.ChangePasswordInfo(db, userDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a token for the user
	tokenString, err := authenticate.CreateToken([]byte(user.Username))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create response object
	createUserResponse := models.CreateUserResponse{
		Username: user.Username,
		Token:    tokenString,
	}

	// write the response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createUserResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
