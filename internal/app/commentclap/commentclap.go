package commentclap

import (
	"time"

	"github.com/Vatsal-S-Patel/Bloggy/internal/consts"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type service struct {
	DB *sqlx.DB
}

type Service interface {
	Add(commentID, userID uuid.UUID) error
	Remove(commentID, userID uuid.UUID) error
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Add(commentID, userID uuid.UUID) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO clapped_comments (user_id, comment_id, created_at) VALUES ($1, $2, $3)`
	_, err = tx.Exec(query, userID, commentID, time.Now())
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "unique_clapped_comments_user_id_comment_id" {
			return errs.ErrAlreadyClapped
		} else if ok && pqErr.Code == consts.DB_CODE_FOREIGN_KEY_CONSTRAINT_VIOLATION && pqErr.Constraint == "clapped_comments_comment_id_fkey" {
			return errs.ErrCommentNotFound
		}
		return err
	}

	query = `UPDATE comments SET claps=claps+1 WHERE id=$1`
	_, err = tx.Exec(query, commentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) Remove(commentID, userID uuid.UUID) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `DELETE FROM clapped_comments WHERE user_id=$1 AND comment_id=$2`
	res, err := s.DB.Exec(query, userID, commentID)
	if err != nil {
		return err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrAlreadyUnClapped
	}

	query = `UPDATE comments SET claps=claps-1 WHERE id=$1`
	_, err = s.DB.Exec(query, commentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
