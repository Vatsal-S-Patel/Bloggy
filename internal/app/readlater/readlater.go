package readlater

import (
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
	Add(readLaterBlog *models.ReadLater) error
	Remove(userID, blogID uuid.UUID) error
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Add(readLaterBlog *models.ReadLater) error {
	query := `INSERT INTO read_later (user_id, blog_id, created_at) VALUES ($1, $2, $3)`

	_, err := s.DB.Exec(query, readLaterBlog.UserID, readLaterBlog.BlogID, readLaterBlog.CreatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "unique_read_later_user_id_blog_id" {
			return errs.ErrAlreadyInReadLater
		}
		return err
	}

	return nil
}

func (s *service) Remove(userID, blogID uuid.UUID) error {
	query := `DELETE FROM read_later WHERE user_id=$1 AND blog_id=$2`

	res, err := s.DB.Exec(query, userID, blogID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrBlogNotFound
	}

	return nil
}
