package comment

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
	Add(comment *models.Comment) error
	AddReply(comment *models.Comment) error
	Get(blogID uuid.UUID) ([]*models.Comment, error)
	GetReplies(parentID uuid.UUID) ([]*models.Comment, error)
	Update(comment *models.Comment) error
	Remove(commentID uuid.UUID) error
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Add(comment *models.Comment) error {
	query := `INSERT INTO comments (id, body, claps, replies, author_id, blog_id, parent_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := s.DB.Exec(query, comment.ID, comment.Body, comment.Claps, comment.Replies, comment.AuthorID, comment.BlogID, nil, comment.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddReply(comment *models.Comment) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO comments (id, body, claps, replies, author_id, blog_id, parent_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(query, comment.ID, comment.Body, comment.Claps, comment.Replies, comment.AuthorID, comment.BlogID, comment.ParentID, comment.CreatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_FOREIGN_KEY_CONSTRAINT_VIOLATION && pqErr.Constraint == "comments_parent_id_fkey" {
			return errs.ErrCommentNotFound
		}
		return err
	}

	query = `UPDATE comments SET replies=replies+1 WHERE id=$1`
	_, err = tx.Exec(query, comment.ParentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) Get(blogID uuid.UUID) ([]*models.Comment, error) {
	query := `SELECT id, body, claps, replies, author_id, blog_id, created_at FROM comments WHERE blog_id=$1 AND parent_id IS NULL ORDER BY created_at DESC`

	var comments []*models.Comment
	err := s.DB.Select(&comments, query, blogID)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errs.ErrCommentNotFound
	}

	return comments, nil
}

func (s *service) GetReplies(parentID uuid.UUID) ([]*models.Comment, error) {
	query := `SELECT id, body, claps, replies, author_id, blog_id, parent_id, created_at FROM comments WHERE parent_id=$1 ORDER BY created_at DESC`

	var comments []*models.Comment
	err := s.DB.Select(&comments, query, parentID)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errs.ErrCommentNotFound
	}

	return comments, nil
}

func (s *service) Update(comment *models.Comment) error {
	query := `UPDATE comments SET body=$1 WHERE id=$2 AND author_id=$3`

	res, err := s.DB.Exec(query, comment.Body, comment.ID, comment.AuthorID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrCommentNotFound
	}

	return nil
}

func (s *service) Remove(commentID uuid.UUID) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `DELETE FROM comments WHERE id=$1 RETURNING parent_id`

	var parentID uuid.UUID
	err = tx.QueryRow(query, commentID).Scan(&parentID)
	if err != nil {
		return err
	}

	query = `UPDATE comments SET replies=replies-1 WHERE id=$1`
	res, err := tx.Exec(query, parentID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrCommentNotFound
	}

	return tx.Commit()
}
