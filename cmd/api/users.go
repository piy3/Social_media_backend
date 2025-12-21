package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/piy3/social/internal/store"
)

type userCtxKey string

const userKey userCtxKey = "user"

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type UpdateUserPayload struct {
	Username *string `json:"username" validate:"omitempty,alphanum,min=3,max=30"`
	Email    *string `json:"email" validate:"omitempty,email"`
}

// func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload CreateUserPayload
// 	if err := readJSON(r, &payload); err != nil {
// 		writeJSONError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	if err := Validate.Struct(payload); err != nil {
// 		writeJSONError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	user := &store.User{
// 		Username: payload.Username,
// 		Email:    payload.Email,
// 		Password: payload.Password,
// 	}
// 	if err := app.store.Users.Create(r.Context(), user); err != nil {
// 		writeJSONError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	if err := writeJSON(w, http.StatusCreated, user); err != nil {
// 		writeJSONError(w, http.StatusInternalServerError, err.Error())
// 	}
// }

// func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	token := chi.URLParam(r, "token")
// 	err := app.store.Users.Activate(r.Context(), token)
// 	if err != nil {
// 		writeJSONError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// }

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var payload UpdateUserPayload
	if err := readJSON(r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Username != nil {
		user.Username = *payload.Username
	}
	if payload.Email != nil {
		user.Email = *payload.Email
	}

	ctx := r.Context()
	if err := app.store.Users.Update(ctx, user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	ctx := r.Context()
	if err := app.store.Users.Delete(ctx, user.ID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userKey).(*store.User)
	return user
}
