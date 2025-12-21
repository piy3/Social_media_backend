package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/piy3/social/internal/store"
)

func Seed(store store.Storage, db *sql.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx,nil)
	if err != nil {
		return err
	}
	// users
	users := generateUser(100)
	for _, u := range users {
		err := store.Users.Create(ctx,tx, u)
		if err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return err
		}
	}
	tx.Commit()
	// posts
	posts := generatePosts(200, users)
	for _, p := range posts {
		err := store.Posts.Create(ctx,p)
		if err != nil {
			log.Println("Error creating post:", err)
			return err
		}
	}

	comments := generateComments(500, posts, users)
	for _, c := range comments {
		err := store.Comments.Create(ctx, c)
		if err != nil {
			log.Println("Error creating comment:", err)
			return err
		}
	}
	log.Println("Seeding completed successfully")
	return nil
}

func generateUser(n int) []*store.User {
	user := make([]*store.User, n)
	for i := 0; i < n; i++ {
		user[i] = &store.User{
			Username: "user" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@example.com",
			// Password: "password123",
		}
	}
	return user
}
func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for i := 0; i < n; i++ {
		posts[i] = &store.Post{
			Title:   "Post Title " + strconv.Itoa(i),
			Content: "This is the content of post " + strconv.Itoa(i),
			UserID:  users[i%len(users)].ID,
			Tags:    []string{"tag1", "tag2"},
		}
	}
	return posts
}

func generateComments(n int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, n)
	for i := 0; i < n; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[i%len(posts)].ID,
			UserID:  users[i%len(users)].ID,
			Content: "This is a comment " + strconv.Itoa(i),
		}
	}
	return comments
}
