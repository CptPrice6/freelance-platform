package validators

import (
	"backend/types"
	"encoding/json"
	"fmt"
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
	}

	return registerRequest, nil

}
