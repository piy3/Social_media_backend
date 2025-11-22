package store

import (
	"context"
	"database/sql"
)

//this store implements repository pattern to interact with data storage

type Storage struct{	
	Posts interface{
		GetByID(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
	}
	Users interface{
		Create(context.Context, *User) error
	}
	Comments interface{
		GetByPostId(context.Context, int64) ([]*Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:&PostStore{db},
		Users:&UserStore{db},
		Comments:&CommentStore{db},
	}
}