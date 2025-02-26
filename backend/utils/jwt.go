package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-salt")

type Claims struct {
	Id        int    `json:"id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token
func GenerateAccessToken(id int, role string) (string, error) {
	claims := Claims{
		Id:        id,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // 30 minutes expiry
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(id int) (string, error) {
	claims := Claims{
		Id:        id,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 day expiry
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateTokenPair(id int, role string) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(id, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateRefreshToken(id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ParseJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims, err := ParseJWT(tokenStr)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type: expected access token")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := ParseJWT(tokenStr)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type: expected refresh token")
	}

	return claims, nil
}
