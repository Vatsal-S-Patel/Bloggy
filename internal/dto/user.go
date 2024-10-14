package dto

type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required,max=30"`
	Email    string `json:"email" validate:"required,email,max=70"`
	Password string `json:"password" validate:"required,min=8,password"`
	Bio      string `json:"bio" validate:"max=500"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
}

type AdminRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type SuperAdminRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FollowRequest struct {
	FollowerID  string `json:"follower_id"`
	FollowingID string `json:"following_id"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,max=30"`
	Password string `json:"password" validate:"required,min=8,password"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}
