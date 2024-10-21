package tag

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
	Add(tag *models.Tag) error
	Get(tagIdentifier string) ([]*models.Tag, error)
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
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "idx_tags_name" {
			return errs.ErrTagAlreadyInUse
		}
		return err
	}

	return nil
}

func (s *service) Get(tagIdentifier string) ([]*models.Tag, error) {
	var tags []*models.Tag

	query := `SELECT id, name FROM tags WHERE LOWER(name) LIKE $1 || '%'`

	tagID, _ := uuid.Parse(tagIdentifier)
	if tagID != uuid.Nil {
		query = `SELECT id, name FROM tags WHERE id=$1`
	}

	err := s.DB.Select(&tags, query, tagIdentifier)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, errs.ErrTagNotFound
	}

	return tags, nil
}
