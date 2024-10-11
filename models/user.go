package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Bio         string    `json:"bio"`
	Avatar      string    `json:"avatar"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	JoinedAt    time.Time `json:"joined_at"`
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
}

type Admin struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Avatar      string    `json:"avatar"`
	JoinedAt    time.Time `json:"joined_at"`
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
}

type SuperAdmin struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	JoinedAt    time.Time `json:"joined_at"`
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
}

type Follow struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Wishlist struct {
	UserID    uuid.UUID `json:"user_id"`
	BlogID    uuid.UUID `json:"blog_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ReadingHistory struct {
	UserID    uuid.UUID `json:"user_id"`
	BlogID    uuid.UUID `json:"blog_id"`
	CreatedAt time.Time `json:"created_at"`
}
