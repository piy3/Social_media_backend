package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/piy3/social/internal/store"
)

type registerUserPayload struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	//reading json from request
	err := readJSON(r, &payload)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	//validating payload
	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

	//creating user
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	//hash - password
	if err:= user.Password.Set(payload.Password);err!=nil{
		writeJSONError(w,http.StatusInternalServerError,err.Error())
		return
	}
	//encrypt token
	plainToken:=uuid.New().String()
	hash:=sha256.Sum256([]byte(plainToken))
	hashToken:= hex.EncodeToString(hash[:])
	//store the user
	if err := app.store.Users.CreateAndInvite(r.Context(), user,hashToken,app.config.mail.exp); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//mail
	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
