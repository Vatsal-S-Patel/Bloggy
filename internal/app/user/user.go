package user

import (
	"database/sql"
	"errors"

	"github.com/Vatsal-S-Patel/Bloggy/internal/consts"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type service struct {
	DB *sqlx.DB
}

type Service interface {
	RegisterUser(user *models.User) error
	GetIDPasswordByUsername(username string) (uuid.UUID, string, error)
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) RegisterUser(user *models.User) error {
	query := `INSERT INTO users (id, username, email, password, bio, avatar, followers, following, joined_at, last_login_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := s.DB.Exec(query, user.ID, user.Username, user.Email, user.Password, user.Bio, user.Avatar, user.Followers, user.Following, user.JoinedAt, user.LastLoginAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION {
			switch pqErr.Constraint {
			case "users_username_key":
				return errs.ErrUsernameAlreadyInUse
			case "users_email_key":
				return errs.ErrUserEmailAlreadyInUse
			}
		}
		return err
	}

	return nil
}

func (s *service) GetIDPasswordByUsername(username string) (uuid.UUID, string, error) {
	query := `SELECT id, password FROM users WHERE username=$1`

	var userID uuid.UUID
	var password string

	row := s.DB.QueryRow(query, username)

	err := row.Scan(&userID, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", errs.ErrUserNotFound
		}
		return uuid.Nil, "", err
	}

	return userID, password, nil
}
