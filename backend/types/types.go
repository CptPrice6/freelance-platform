package types

// request structure for register
type RegisterRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"Surname"`
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

type UpdateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	Name        string `json:"name"`
	Surname     string `json:"Surname"`
}

type UpdateUserRequestAdmin struct {
	Role string `json:"role"`
	Ban  bool   `json:"ban"`
}

type UpdateFreelancerDataRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	HourlyRate   float64 `json:"hourly_rate"`
	WorkType     string  `json:"work_type"`
	HoursPerWeek int     `json:"hours_per_week"`
}

type UpdateClientDataRequest struct {
	Description string `json:"description"`
	CompanyName string `json:"company_name"`
	Industry    string `json:"industry"`
	Location    string `json:"location"`
}
type AddDeleteSkillRequest struct {
	SkillID int `json:"skill_id"`
}

type AddUpdateSkillRequest struct {
	SkillName string `json:"skill_name"`
}
