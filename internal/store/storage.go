package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("resource not found")

//this store implements repository pattern to interact with data storage

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
	}		
	Users interface {
		Create(context.Context,*sql.Tx, *User) error
		CreateAndInvite(ctx context.Context,user *User,token string,invitationExpiry time.Duration) error
		GetByID(context.Context, int64) (*User, error)
		Update(context.Context, *User) error
		Delete(context.Context, int64) error
		Activate(context.Context, string) error
	}	
	Comments interface {	
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, int64) ([]*Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db}, 
	} 
}


func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx,err:= db.BeginTx(ctx,nil)
	if err!=nil{
		return err
	}
	if err:= fn(tx);err!=nil{
		if rbErr:= tx.Rollback();rbErr!=nil{
			return rbErr
		}
		return err
	}
	return tx.Commit()
}