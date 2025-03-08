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

func LogoutUserValidator(requestBody []byte) (*types.LogoutUserRequest, error) {

	var logoutUserRequest = new(types.LogoutUserRequest)

	err := json.Unmarshal(requestBody, &logoutUserRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if logoutUserRequest.UserId == 0 {
		return nil, fmt.Errorf("Missing required fields: user_id")
	}

	return logoutUserRequest, nil

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
	} else if updateUserRequest.NewPassword != "" {
		err = ValidatePassword(updateUserRequest.NewPassword)
		if err != nil {
			return nil, err
		}
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
