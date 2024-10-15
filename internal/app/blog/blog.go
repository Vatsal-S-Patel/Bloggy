package blog

import (
	"strings"

	"github.com/Vatsal-S-Patel/Bloggy/internal/consts"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type service struct {
	DB *sqlx.DB
}

type Service interface {
	Publish(*models.Blog, []*models.BlogTag) error
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

	query := `INSERT INTO blogs (id, title, subtitle, content, ft_image, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(query, blog.ID, blog.Title, blog.Subtitle, blog.Content, blog.FtImage, blog.AuthorID, blog.CreatedAt, blog.UpdatedAt)
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok && pqErr.Code == consts.DB_CODE_FOREIGN_KEY_CONSTRAINT_VIOLATION && pqErr.Constraint == "blog_tags_tag_id_fkey" {
				return errs.ErrTagNotFound
			}
			return err
		}
	}

	tx.Commit()
	return nil
}
