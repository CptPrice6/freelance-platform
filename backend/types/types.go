package types

// request structure for register
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// request structure for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutUserRequest struct {
	UserId int `json:"user_id"`
}

type UpdateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
