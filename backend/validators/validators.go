package validators

import (
	"backend/types"
	"encoding/json"
	"fmt"
)

func RegisterValidator(requestBody []byte) (*types.RegisterLoginRequest, error) {

	var registerRequest = new(types.RegisterLoginRequest)

	err := json.Unmarshal(requestBody, &registerRequest)
	if err != nil {
		fmt.Println("Error parsing request body:", err)
		return nil, fmt.Errorf("Invalid input")
	}

	if registerRequest.Email == "" {
		return nil, fmt.Errorf("Missing required fields: email")
	} else if registerRequest.Password == "" {
		return nil, fmt.Errorf("Missing required fields: password")
	}

	// add a check for password length and email validity

	return registerRequest, nil

}

func LoginValidator(requestBody []byte) (*types.RegisterLoginRequest, error) {

	var loginRequest = new(types.RegisterLoginRequest)

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
