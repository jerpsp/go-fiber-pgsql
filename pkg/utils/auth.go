package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrAccessDenied       = errors.New("access denied: insufficient permissions")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

// ValidateToken validates a JWT token and returns the user information
func ValidateToken(cfg *config.Config, tokenString string, tokenType TokenType) (*UserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Verify token type
	if claims["type"].(string) != string(tokenType) {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// Parse user ID
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get role if exists, default to "user" if not present (for backwards compatibility)
	role := "user"
	if roleValue, exists := claims["role"]; exists {
		if roleStr, ok := roleValue.(string); ok {
			role = roleStr
		}
	}

	return &UserInfo{
		ID:    userID,
		Email: claims["email"].(string),
		Role:  role,
	}, nil
}
