package blog

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Vatsal-S-Patel/Bloggy/internal/consts"
	"github.com/Vatsal-S-Patel/Bloggy/internal/dto"
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
	Publish(blog *models.Blog, blogTags []*models.BlogTag) error
	Get(blogID, userID uuid.UUID) (*dto.Blog, error)
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Publish(blog *models.Blog, blogTags []*models.BlogTag) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}

	query := `SELECT username FROM users WHERE id=$1`

	var author string
	err = tx.Get(&author, query, blog.AuthorID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	query = `INSERT INTO blogs (id, title, subtitle, content, ft_image, author_id, author, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = tx.Exec(query, blog.ID, blog.Title, blog.Subtitle, blog.Content, blog.FtImage, blog.AuthorID, author, blog.CreatedAt, blog.UpdatedAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if len(blogTags) > 0 {
		queryBuilder := strings.Builder{}
		queryBuilder.WriteString(`INSERT INTO blog_tags (blog_id, tag_id) VALUES `)

		for _, blogTag := range blogTags {
			queryBuilder.WriteString(`('` + blogTag.BlogID.String() + `','` + blogTag.TagID.String() + `'),`)
		}

		_, err = tx.Exec(queryBuilder.String()[:queryBuilder.Len()-1])
		if err != nil {
			_ = tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok && pqErr.Code == consts.DB_CODE_FOREIGN_KEY_CONSTRAINT_VIOLATION && pqErr.Constraint == "blog_tags_tag_id_fkey" {
				return errs.ErrTagNotFound
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Get(blogID, userID uuid.UUID) (*dto.Blog, error) {
	readHistoryErrChan := make(chan error, 1)

	if userID != uuid.Nil {
		go func() {
			// insert into reading_history but if already there then update time
			query := `INSERT INTO reading_history (user_id, blog_id, created_at) VALUES ($1, $2, $3) 
			ON CONFLICT (user_id, blog_id) 
			DO UPDATE SET created_at=$3`
			_, err := s.DB.Exec(query, userID, blogID, time.Now())
			if err != nil {
				readHistoryErrChan <- err
			}
			readHistoryErrChan <- nil
		}()
	} else {
		readHistoryErrChan <- nil
	}

	query := `SELECT id, title, subtitle, content, ft_image, author_id, author, created_at, updated_at FROM blogs WHERE id=$1`

	var blog dto.Blog
	err := s.DB.Get(&blog, query, blogID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBlogNotFound
		}
		return nil, err
	}

	query = `SELECT t.id, t.name FROM blog_tags as bt INNER JOIN tags as t ON bt.tag_id = t.id WHERE bt.blog_id=$1`

	var tags []*models.Tag
	// not handling error because it's not critical
	_ = s.DB.Select(&tags, query, blogID)

	blog.Tags = tags

	err = <-readHistoryErrChan
	if err != nil {
		return nil, err
	}

	return &blog, nil
}
