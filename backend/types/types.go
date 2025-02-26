package types

// request structure for register and login
type RegisterLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutUserRequest struct {
	UserId int `json:"user_id"`
}
