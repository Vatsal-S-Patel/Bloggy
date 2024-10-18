package tag

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
	Add(tag *models.Tag) error
	Get(tagIdentifier string) (*models.Tag, error)
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Add(tag *models.Tag) error {
	query := `INSERT INTO tags (id, name) VALUES ($1, $2)`

	_, err := s.DB.Exec(query, tag.ID, tag.Name)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "tags_name_key" {
			return errs.ErrTagAlreadyInUse
		}
		return err
	}

	return nil
}

func (s *service) Get(tagIdentifier string) (*models.Tag, error) {
	query := `SELECT id, name FROM tags WHERE name=$1`

	tagID, _ := uuid.Parse(tagIdentifier)
	if tagID != uuid.Nil {
		query = `SELECT id, name FROM tags WHERE id=$1`
	}

	var tag models.Tag
	err := s.DB.Get(&tag, query, tagIdentifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTagNotFound
		}
		return nil, err
	}

	return &tag, nil
}
