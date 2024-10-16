package history

import (
	"github.com/Vatsal-S-Patel/Bloggy/internal/dto"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type service struct {
	DB *sqlx.DB
}

type Service interface {
	Get(userID uuid.UUID) ([]*dto.History, error)
	Remove(userID, blogID uuid.UUID) error
	RemoveAll(userID uuid.UUID) error
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Get(userID uuid.UUID) ([]*dto.History, error) {
	query := `SELECT b.id, b.title, b.author_id, b.author, b.ft_image, rh.created_at 
	FROM reading_history AS rh INNER JOIN blogs AS b ON rh.blog_id=b.id 
	WHERE rh.user_id=$1 
	ORDER BY created_at DESC`

	var history []*dto.History
	err := s.DB.Select(&history, query, userID)
	if err != nil {
		return nil, err
	}
	if len(history) == 0 {
		return nil, errs.ErrHistoryNotFound
	}

	return history, nil
}

func (s *service) Remove(userID, blogID uuid.UUID) error {
	query := `DELETE FROM reading_history WHERE user_id=$1 AND blog_id=$2`

	res, err := s.DB.Exec(query, userID, blogID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrHistoryNotFound
	}

	return nil
}

func (s *service) RemoveAll(userID uuid.UUID) error {
	query := `DELETE FROM reading_history WHERE user_id=$1`

	res, err := s.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrHistoryNotFound
	}

	return nil
}
