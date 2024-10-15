package errs

import "errors"

var (
	ErrUserEmailAlreadyInUse = errors.New("email already in use")
	ErrUsernameAlreadyInUse  = errors.New("username already in use")
	ErrUserNotFound          = errors.New("user not found")

	ErrBlogNotFound = errors.New("blog not found")

	ErrTagNotFound = errors.New("tag not found")

	ErrDraftNotFound = errors.New("drafts not found")
)
