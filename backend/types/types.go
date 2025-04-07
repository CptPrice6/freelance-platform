package types

import "time"

// request structure for register
type RegisterRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
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
	Surname     string `json:"surname"`
}

type UpdateUserRequestAdmin struct {
	Role string `json:"role"`
	Ban  bool   `json:"ban"`
}

type UpdateFreelancerDataRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	HourlyRate   float64 `json:"hourly_rate"`
	HoursPerWeek string  `json:"hours_per_week"`
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

type Skill struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FreelancerData struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Skills       []Skill `json:"skills"`
	HourlyRate   float64 `json:"hourly_rate"`
	HoursPerWeek string  `json:"hours_per_week"`
}

type ClientData struct {
	Description string `json:"description"`
	CompanyName string `json:"company_name"`
	Industry    string `json:"industry"`
	Location    string `json:"location"`
}

type UserResponse struct {
	ID             int             `json:"id"`
	Email          string          `json:"email"`
	Role           string          `json:"role"`
	Name           string          `json:"name"`
	Surname        string          `json:"surname"`
	FreelancerData *FreelancerData `json:"freelancer_data,omitempty"`
	ClientData     *ClientData     `json:"client_data,omitempty"`
}

type UserResponseForAdmins struct {
	ID             int             `json:"id"`
	Email          string          `json:"email"`
	Role           string          `json:"role"`
	Name           string          `json:"name"`
	Surname        string          `json:"surname"`
	Ban            bool            `json:"ban"`
	CreatedAt      time.Time       `json:"created_at"`
	FreelancerData *FreelancerData `json:"freelancer_data,omitempty"`
	ClientData     *ClientData     `json:"client_data,omitempty"`
}

type FreelancerInfo struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Title      string  `json:"title"`
	HourlyRate float64 `json:"hourly_rate"`
}

type ClientInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	CompanyName string `json:"company_name"`
}

type CreateJobRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Rate         string   `json:"rate"`
	Amount       int      `json:"amount"`
	Length       string   `json:"length"`
	HoursPerWeek string   `json:"hours_per_week"`
	Skills       []*Skill `json:"skills"`
}

type UpdateJobRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Rate         string   `json:"rate"`
	Amount       int      `json:"amount"`
	Length       string   `json:"length"`
	HoursPerWeek string   `json:"hours_per_week"`
	Skills       []*Skill `json:"skills"`
}

type JobInfo struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Type          string  `json:"type"`
	Rate          string  `json:"rate"`
	Amount        int     `json:"amount"`
	Length        string  `json:"length"`
	HoursPerWeek  string  `json:"hours_per_week"`
	ClientID      int     `json:"client_id"`
	Skills        []Skill `json:"skills"`
	ApplicationID int     `json:"application_id"`
}

type ClientJobInfo struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Type             string  `json:"type"`
	Rate             string  `json:"rate"`
	Amount           int     `json:"amount"`
	Length           string  `json:"length"`
	HoursPerWeek     string  `json:"hours_per_week"`
	Status           string  `json:"status"`
	ClientID         int     `json:"client_id"`
	FreelancerID     int     `json:"freelancer_id"`
	Skills           []Skill `json:"skills"`
	ApplicationCount int     `json:"application_count"`
}

type FreelancerJobInfo struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Type          string  `json:"type"`
	Rate          string  `json:"rate"`
	Amount        int     `json:"amount"`
	Length        string  `json:"length"`
	HoursPerWeek  string  `json:"hours_per_week"`
	Status        string  `json:"status"`
	ClientID      int     `json:"client_id"`
	Skills        []Skill `json:"skills"`
	ApplicationID int     `json:"application_id"`
}

type ClientJobDetailedInfo struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Type         string        `json:"type"`
	Rate         string        `json:"rate"`
	Amount       int           `json:"amount"`
	Length       string        `json:"length"`
	HoursPerWeek string        `json:"hours_per_week"`
	Status       string        `json:"status"`
	ClientID     int           `json:"client_id"`
	FreelancerID int           `json:"freelancer_id"`
	Skills       []Skill       `json:"skills"`
	Applications []Application `json:"applications"`
}

type Application struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	JobID           int         `json:"job_id"`
	Description     string      `json:"description"`
	RejectionReason string      `json:"rejection_reason"`
	Status          string      `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	Attachment      *Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	ID            int       `json:"id"`
	ApplicationID int       `json:"application_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	CreatedAt     time.Time `json:"created_at"`
}

type SubmitApplicationRequest struct {
	JobID       int    `json:"job_id"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
	FileBase64  string `json:"file_base64"`
}

type UpdateApplicationRequest struct {
	Description string `json:"description"`
	FileName    string `json:"file_name"`
	FileBase64  string `json:"file_base64"`
}

type ChangeApplicationStatusRequest struct {
	Status          string `json:"status"`
	RejectionReason string `json:"rejection_reason"`
}
