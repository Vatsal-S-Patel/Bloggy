package draft

import (
	"database/sql"
	"errors"

	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type service struct {
	DB *sqlx.DB
}

type Service interface {
	Add(*models.Draft) error
	Get(uuid.UUID) (*models.Draft, error)
	GetAll(uuid.UUID) ([]*models.Draft, error)
	Update(*models.Draft) error
	Remove(id uuid.UUID) error
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Add(draft *models.Draft) error {
	query := `INSERT INTO drafts (id, title, subtitle, content, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.DB.Exec(query, draft.ID, draft.Title, draft.Subtitle, draft.Content, draft.AuthorID, draft.CreatedAt, draft.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAll(authorID uuid.UUID) ([]*models.Draft, error) {
	query := `SELECT id, title, subtitle, content, author_id, created_at, updated_at FROM drafts WHERE author_id=$1`

	var drafts []*models.Draft
	err := s.DB.Select(&drafts, query, authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDraftNotFound
		}
		return nil, err
	}

	return drafts, nil
}

func (s *service) Get(draftID uuid.UUID) (*models.Draft, error) {
	query := `SELECT id, title, subtitle, content, author_id, created_at, updated_at FROM drafts WHERE id=$1`

	var draft models.Draft
	err := s.DB.Get(&draft, query, draftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDraftNotFound
		}
		return nil, err
	}

	return &draft, nil
}

func (s *service) Update(draft *models.Draft) error {
	query := `UPDATE drafts SET title=$1, subtitle=$2, content=$3, updated_at=$4 WHERE id=$5`

	res, err := s.DB.Exec(query, draft.Title, draft.Subtitle, draft.Content, draft.UpdatedAt, draft.ID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrDraftNotFound
	}

	return nil
}

func (s *service) Remove(draftID uuid.UUID) error {
	query := `DELETE FROM drafts WHERE id=$1`

	res, err := s.DB.Exec(query, draftID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrDraftNotFound
	}

	return nil
}
