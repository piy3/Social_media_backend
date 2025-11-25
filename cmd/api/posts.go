package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/piy3/social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=255"`
	Content string   `json:"content" validate:"required,max=5000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a post
	var payload CreatePostPayload
	err := readJSON(r, &payload)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// UserID and Tags can be set here as needed
		UserID: 1, // Example user ID
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting a post by ID
	// idParam := chi.URLParam(r, "postID")
	// id, err := strconv.ParseInt(idParam, 10, 64)
	// if err != nil {
	// 	writeJSONError(w, http.StatusBadRequest, "invalid post ID")
	// 	return
	// }
	// ctx := r.Context()
	// post, err := app.store.Posts.GetByID(ctx, id)
	// if err != nil {
	// 	writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// if err := writeJSON(w, http.StatusOK, post); err != nil {
	// 	writeJSONError(w, http.StatusInternalServerError, err.Error())
	// }

	//implementation - 2 using middleware to get post from context
	post := getPostFromCtx(r)
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

}

type UpdatedPostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=255"`
	Content *string   `json:"content" validate:"omitempty,max=5000"`
	Tags    *[]string `json:"tags"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a post by ID
	post := getPostFromCtx(r)
	var payload UpdatedPostPayload
	if err := readJSON(r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}
	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting a post by ID
	idParam := chi.URLParam(r, "postID")
	log.Println("Deleting post with ID:", idParam)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid post ID")
		return
	}
	ctx := r.Context()
	err = app.store.Posts.Delete(ctx, id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid post ID")
			return
		}
		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, ok := r.Context().Value(postCtx).(*store.Post)
	if !ok {
		return nil
	}
	return post
}
