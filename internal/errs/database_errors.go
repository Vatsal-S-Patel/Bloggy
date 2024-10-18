package errs

import "errors"

var (
	ErrUserEmailAlreadyInUse = errors.New("email already in use")
	ErrUsernameAlreadyInUse  = errors.New("username already in use")
	ErrUserNotFound          = errors.New("user not found")

	ErrBlogNotFound = errors.New("blog not found")

	ErrTagNotFound = errors.New("tag not found")

	ErrDraftNotFound = errors.New("drafts not found")

	ErrHistoryNotFound = errors.New("history not found")

	ErrAlreadyInReadLater = errors.New("blog already in readlater")

	ErrBookmarkNameAlreadyInUse = errors.New("bookmark name in use")
	ErrBlogAlreadyInBookmark    = errors.New("blog already in bookmark")
	ErrBookmarkNotFound         = errors.New("bookmark not found")
)
