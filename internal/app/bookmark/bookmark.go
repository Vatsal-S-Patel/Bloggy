package bookmark

import (
	"database/sql"
	"errors"
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
	Add(bookmark *models.Bookmark) error
	AddBlog(bookmarkID, blogID uuid.UUID) error
	GetAll(userID uuid.UUID) ([]*models.Bookmark, error)
	Get(bookmarkID uuid.UUID) (*models.Bookmark, error)
	GetBlogs(bookmarkID uuid.UUID) ([]*dto.BookmarkBlogs, error)
	Update(updatedBookmark *dto.AddBookmarkRequest, bookmarkID, userID uuid.UUID) error
	Remove(bookmarkID, userID uuid.UUID) error
	RemoveBlog(bookmarkID, blogID uuid.UUID) error
}

func NewService(db *sqlx.DB) Service {
	return &service{
		DB: db,
	}
}

func (s *service) Add(bookmark *models.Bookmark) error {
	query := `INSERT INTO bookmarks (id, name, user_id, visible, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.DB.Exec(query, bookmark.ID, bookmark.Name, bookmark.UserID, bookmark.Visible, bookmark.CreatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "unique_bookmarks_name" {
			return errs.ErrBookmarkNameAlreadyInUse
		}
		return err
	}

	return nil
}

func (s *service) AddBlog(bookmarkID, blogID uuid.UUID) error {
	query := `INSERT INTO bookmark_blogs (bookmark_id, blog_id, created_at) VALUES ($1, $2, $3)`

	_, err := s.DB.Exec(query, bookmarkID, blogID, time.Now())
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "unique_bookmark_blogs_bookmark_id_blog_id" {
			return errs.ErrBlogAlreadyInBookmark
		} else if ok && pqErr.Code == consts.DB_CODE_FOREIGN_KEY_CONSTRAINT_VIOLATION && pqErr.Constraint == "bookmark_blogs_blog_id_fkey" {
			return errs.ErrBlogNotFound
		}
		return err
	}

	return nil
}

func (s *service) GetAll(userID uuid.UUID) ([]*models.Bookmark, error) {
	query := `SELECT id, name, user_id, created_at, visible FROM bookmarks WHERE user_id=$1 ORDER BY created_at DESC`

	var bookmarks []*models.Bookmark
	err := s.DB.Select(&bookmarks, query, userID)
	if err != nil {
		return nil, err
	}
	if len(bookmarks) == 0 {
		return nil, errs.ErrBookmarkNotFound
	}

	return bookmarks, nil
}

func (s *service) Get(bookmarkID uuid.UUID) (*models.Bookmark, error) {
	query := `SELECT id, name, user_id, created_at, visible FROM bookmarks WHERE id=$1`

	var bookmark models.Bookmark
	err := s.DB.Get(&bookmark, query, bookmarkID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBookmarkNotFound
		}
		return nil, err
	}

	return &bookmark, nil
}

func (s *service) GetBlogs(bookmarkID uuid.UUID) ([]*dto.BookmarkBlogs, error) {
	query := `SELECT b.id, b.title, b.ft_image, b.author_id, b.author, bb.created_at FROM bookmark_blogs AS bb INNER JOIN blogs AS b ON bb.blog_id=b.id WHERE bb.bookmark_id=$1 ORDER BY bb.created_at DESC`

	var blogs []*dto.BookmarkBlogs
	err := s.DB.Select(&blogs, query, bookmarkID)
	if err != nil {
		return nil, err
	}
	if len(blogs) == 0 {
		return nil, errs.ErrBlogNotFound
	}

	return blogs, nil
}

func (s *service) Update(updatedBookmark *dto.AddBookmarkRequest, bookmarkID, userID uuid.UUID) error {
	query := `UPDATE bookmarks SET name=$1, visible=$2 WHERE id=$3 AND user_id=$4`

	res, err := s.DB.Exec(query, updatedBookmark.Name, updatedBookmark.Visible, bookmarkID, userID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == consts.DB_CODE_UNIQUE_CONSTRAINT_VIOLATION && pqErr.Constraint == "unique_bookmarks_name" {
			return errs.ErrBookmarkNameAlreadyInUse
		}
		return err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrBookmarkNotFound
	}

	return nil
}

func (s *service) Remove(bookmarkID, userID uuid.UUID) error {
	query := `DELETE FROM bookmarks WHERE id=$1 AND user_id=$2`

	res, err := s.DB.Exec(query, bookmarkID, userID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errs.ErrBookmarkNotFound
	}

	return nil
}

func (s *service) RemoveBlog(bookmarkID, blogID uuid.UUID) error {
	query := `DELETE FROM bookmark_blogs WHERE bookmark_id=$1 AND blog_id=$2`

	res, err := s.DB.Exec(query, bookmarkID, blogID)
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
