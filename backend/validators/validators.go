package validators

import (
	"backend/types"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-passwd/validator"
)

func RegisterValidator(requestBody []byte) (*types.RegisterRequest, error) {

	var registerRequest = new(types.RegisterRequest)

	err := json.Unmarshal(requestBody, &registerRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if registerRequest.Email == "" {
		return nil, fmt.Errorf("Missing required fields: email")
	} else if registerRequest.Password == "" {
		return nil, fmt.Errorf("Missing required fields: password")
	} else if registerRequest.Role == "" {
		return nil, fmt.Errorf("Missing required fields: role")
	} else if registerRequest.Name == "" {
		return nil, fmt.Errorf("Missing required fields: name")
	} else if registerRequest.Surname == "" {
		return nil, fmt.Errorf("Missing required fields: surname")
	}

	if registerRequest.Role != "client" && registerRequest.Role != "freelancer" {
		return nil, fmt.Errorf("Invalid role: %s. Role must be either 'client' or 'freelancer'", registerRequest.Role)
	}

	err = ValidateEmail(registerRequest.Email)
	if err != nil {
		return nil, err
	}

	err = ValidatePassword(registerRequest.Password)
	if err != nil {
		return nil, err
	}

	if len(registerRequest.Name) > 30 {
		return nil, fmt.Errorf("Name cannot be longer than 30 symbols")
	}

	if len(registerRequest.Surname) > 30 {
		return nil, fmt.Errorf("Surname cannot be longer than 30 symbols")
	}

	return registerRequest, nil

}

func LoginValidator(requestBody []byte) (*types.LoginRequest, error) {

	var loginRequest = new(types.LoginRequest)

	err := json.Unmarshal(requestBody, &loginRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if loginRequest.Email == "" {
		return nil, fmt.Errorf("Missing required fields: email")
	} else if loginRequest.Password == "" {
		return nil, fmt.Errorf("Missing required fields: password")
	}

	return loginRequest, nil

}

func RefreshValidator(requestBody []byte) (*types.RefreshRequest, error) {

	var refreshRequest = new(types.RefreshRequest)

	err := json.Unmarshal(requestBody, &refreshRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if refreshRequest.RefreshToken == "" {
		return nil, fmt.Errorf("Missing required fields: refresh_token")
	}

	return refreshRequest, nil

}

func UpdateUserValidator(requestBody []byte) (*types.UpdateUserRequest, error) {

	var updateUserRequest = new(types.UpdateUserRequest)

	err := json.Unmarshal(requestBody, &updateUserRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}
	if updateUserRequest.Email != "" {
		err = ValidateEmail(updateUserRequest.Email)
		if err != nil {
			return nil, err
		}
	}
	if updateUserRequest.Password != "" && updateUserRequest.NewPassword == "" {
		return nil, fmt.Errorf("Missing new password")
	}
	if updateUserRequest.NewPassword != "" && updateUserRequest.Password == "" {
		return nil, fmt.Errorf("Missing old password")
	}

	if updateUserRequest.NewPassword != "" {
		err = ValidatePassword(updateUserRequest.NewPassword)
		if err != nil {
			return nil, err
		}
	}
	if updateUserRequest.Name != "" {
		if len(updateUserRequest.Name) > 30 {
			return nil, fmt.Errorf("Name cannot be longer than 30 symbols")
		}
	}
	if updateUserRequest.Surname != "" {
		if len(updateUserRequest.Surname) > 30 {
			return nil, fmt.Errorf("Surname cannot be longer than 30 symbols")
		}
	}

	return updateUserRequest, nil

}

func UpdateUserValidatorAdmin(requestBody []byte) (*types.UpdateUserRequestAdmin, error) {

	var updateUserRequest = new(types.UpdateUserRequestAdmin)

	err := json.Unmarshal(requestBody, &updateUserRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if updateUserRequest.Role != "" && updateUserRequest.Role != "client" && updateUserRequest.Role != "freelancer" && updateUserRequest.Role != "admin" {
		return nil, fmt.Errorf("Invalid role: %s. Role must be 'client', 'freelancer' or 'admin' ", updateUserRequest.Role)
	}

	return updateUserRequest, nil

}

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return errors.New("Invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	// Create a new validator instance
	v := validator.New(validator.MinLength(8, errors.New("Password must contain at least 8 characters")),
		validator.MaxLength(20, errors.New("Password must not exceed 20 characters")),
		validator.CommonPassword(errors.New("Password cannot be common, please choose something unique")),
		validator.ContainsAtLeast("0123456789", 1, errors.New("Password must contain at least one number")),
		validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 1, errors.New("Password must contain at least one uppercase letter")),
		validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 1, errors.New("Password must contain at least one lowercase letter")),
	)
	// Applying validation rules
	err := v.Validate(password)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFreelancerDataValidator(requestBody []byte) (*types.UpdateFreelancerDataRequest, error) {

	var updateFreelancerDataRequest = new(types.UpdateFreelancerDataRequest)

	err := json.Unmarshal(requestBody, &updateFreelancerDataRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if updateFreelancerDataRequest.HourlyRate != 0 {
		if updateFreelancerDataRequest.HourlyRate > 1000 {
			return nil, fmt.Errorf("Hourly Rate cannot be more than 1000")
		}
		if updateFreelancerDataRequest.HourlyRate < 1 {
			return nil, fmt.Errorf("Hourly Rate cannot be less than 1")
		}
	}
	if updateFreelancerDataRequest.WorkType != "" {
		if updateFreelancerDataRequest.WorkType != "remote" && updateFreelancerDataRequest.WorkType != "on-site" && updateFreelancerDataRequest.WorkType != "hybrid" {
			return nil, fmt.Errorf("Invalid Work type: must be 'remote', 'on-site' or hybrid")
		}
	}
	if updateFreelancerDataRequest.HoursPerWeek != 0 {
		if updateFreelancerDataRequest.HoursPerWeek < 1 {
			return nil, fmt.Errorf("Hours per week cannot be less than 1")
		}
		if updateFreelancerDataRequest.HoursPerWeek > 168 {
			return nil, fmt.Errorf("Hours per week cannot be more than 168")
		}
	}
	if updateFreelancerDataRequest.Description != "" {
		if len(updateFreelancerDataRequest.Description) > 1000 {
			return nil, fmt.Errorf("Description cannot be longer than 1000 symbols")
		}
	}
	if updateFreelancerDataRequest.Title != "" {
		if len(updateFreelancerDataRequest.Title) > 30 {
			return nil, fmt.Errorf("Title cannot be longer than 30 symbols")
		}
	}

	return updateFreelancerDataRequest, nil

}

func UpdateClientDataValidator(requestBody []byte) (*types.UpdateClientDataRequest, error) {

	var updateClientDataRequest = new(types.UpdateClientDataRequest)

	err := json.Unmarshal(requestBody, &updateClientDataRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}
	if len(updateClientDataRequest.CompanyName) > 30 {
		return nil, fmt.Errorf("Company name cannot be more than 30 symbols")
	}
	if len(updateClientDataRequest.Industry) > 30 {
		return nil, fmt.Errorf("Industry name cannot be more than 30 symbols")
	}
	if len(updateClientDataRequest.Location) > 30 {
		return nil, fmt.Errorf("Location name cannot be more than 30 symbols")
	}

	return updateClientDataRequest, nil

}

func AddDeleteSkillValidator(requestBody []byte) (*types.AddDeleteSkillRequest, error) {

	var addDeleteSkillRequest = new(types.AddDeleteSkillRequest)

	err := json.Unmarshal(requestBody, &addDeleteSkillRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	return addDeleteSkillRequest, nil

}

func AddUpdateSkillValidator(requestBody []byte) (*types.AddUpdateSkillRequest, error) {

	var addUpdateSkillRequest = new(types.AddUpdateSkillRequest)

	err := json.Unmarshal(requestBody, &addUpdateSkillRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if addUpdateSkillRequest.SkillName == "" {
		return nil, fmt.Errorf("Skill name cannot be empty")
	}

	if len(addUpdateSkillRequest.SkillName) > 50 {
		return nil, fmt.Errorf("Skill name cannot be more than 50 symbols")
	}

	return addUpdateSkillRequest, nil

}

func CreateJobValidator(requestBody []byte) (*types.CreateJobRequest, error) {

	var createJobRequest = new(types.CreateJobRequest)

	err := json.Unmarshal(requestBody, &createJobRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}
	if createJobRequest.Title == "" {
		return nil, fmt.Errorf("Missing required fields: title")
	} else if createJobRequest.Description == "" {
		return nil, fmt.Errorf("Missing required fields: description")
	} else if createJobRequest.Type == "" {
		return nil, fmt.Errorf("Missing required fields: type")
	} else if createJobRequest.Rate == "" {
		return nil, fmt.Errorf("Missing required fields: rate")
	} else if createJobRequest.Amount == 0 {
		return nil, fmt.Errorf("Missing required fields: amount")
	} else if createJobRequest.Length == "" {
		return nil, fmt.Errorf("Missing required fields: length")
	} else if createJobRequest.HoursPerWeek == "" {
		return nil, fmt.Errorf("Missing required fields: hours_per_week")
	}

	if len(createJobRequest.Title) > 30 {
		return nil, fmt.Errorf("Title cannot be longer than 30 symbols")
	}
	if len(createJobRequest.Description) > 1000 {
		return nil, fmt.Errorf("Description cannot be longer than 1000 symbols")
	}

	if !types.ValidProjectTypes[createJobRequest.Type] {
		return nil, errors.New("invalid project type: must be 'ongoing' or 'one-time'")
	}

	if !types.ValidProjectRates[createJobRequest.Rate] {
		return nil, errors.New("invalid project rate: must be 'hourly' or 'fixed'")
	}
	if createJobRequest.Amount < 1 {
		return nil, errors.New("amount cannot be less than 1")
	}
	if createJobRequest.Rate == "hourly" && createJobRequest.Amount > 1000 {
		return nil, errors.New("amount cannot be more than 1000 if rate is hourly")
	}

	if !types.ValidProjectLengths[createJobRequest.Length] {
		return nil, errors.New("invalid project length: must be '<1', '1-3', '3-6', '6-12', or '12+'")
	}

	if !types.ValidProjectHoursPerWeek[createJobRequest.HoursPerWeek] {
		return nil, errors.New("invalid project hours per week: must be '<10', '10-20', '20-40', '40-60', or '80+'")
	}

	return createJobRequest, nil

}
