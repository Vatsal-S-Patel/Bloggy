package dto

type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

type AdminRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type SuperAdminRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FollowRequest struct {
	FollowerID  string `json:"follower_id"`
	FollowingID string `json:"following_id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
