package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	IsActive  bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(plainText string) error {
	// Implementation for hashing and setting the password
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &plainText
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	// Implementation for creating a user in the database
	query := `INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, email,  created_at, updated_at, is_active FROM users WHERE id = $1`
	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		// &user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive, 
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `UPDATE users 
		SET username = COALESCE(NULLIF($1, ''), username),
			email = COALESCE(NULLIF($2, ''), email),
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExpiry time.Duration) error {
	//transaction wrapper
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		//crate the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		//create the user invitation

		err := createInvitation(ctx, tx, user.ID, invitationExpiry, token)
		if err != nil {
			return err
		}
		return nil
	})

}

func createInvitation(ctx context.Context, tx *sql.Tx, userID int64, expiry time.Duration, token string) error {
	query := `INSERT INTO user_invitations (user_id, token, expiry, created_at)
	VALUES ($1,$2, $3,NOW())`
	_, err := tx.ExecContext(
		ctx,
		query,
		userID,
		token,
		time.Now().Add(expiry))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(ctx,s.db,func(tx *sql.Tx) error{
		//1.find the user having this token
	user,err:= s.getUserByInvitationToken(ctx,token)
	if err!=nil{
		return err
	}
	//2.update the user to set activated to true
	user.IsActive= true
	if err:= s.update(ctx,tx,user);err!=nil{
		return err
	}
	//3.delete the token from user_invitations table
	if err:= s.deleteInvitation(ctx,tx,user.ID);err!=nil{
		return err
	}
	return nil
	})
	
}

func (s *UserStore) getUserByInvitationToken(ctx context.Context, token string) (*User, error) {
	query := `SELECT u.id, u.username, u.email,  u.created_at, u.updated_at
	FROM users u JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > NOW()`
	hash:=sha256.Sum256([]byte(token))
	hashToken:= hex.EncodeToString(hash[:])
	user:= &User{}
	err := s.db.QueryRowContext(ctx, query,hashToken).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (s *UserStore) update(ctx context.Context,tx *sql.Tx, user *User) error {
	query := `UPDATE users 
		SET is_active = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at`
	err := tx.QueryRowContext(	
		ctx,
		query,
		user.IsActive,
		user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userID)	
	if err != nil {
		return err
	}
	return nil
}