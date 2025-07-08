package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"

	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type AuthService interface {
	Login(req LoginRequest) (*TokenResponse, error)
	RefreshToken(refreshToken string) (*TokenResponse, error)
	Logout(userID uuid.UUID) error
}

type authService struct {
	config         *config.Config
	userRepository user.UserRepository
	tokenRepo      AuthRepository
}

func NewAuthService(config *config.Config, userRepository user.UserRepository, tokenRepo AuthRepository) AuthService {
	return &authService{config: config, userRepository: userRepository, tokenRepo: tokenRepo}
}

func (s *authService) Login(req LoginRequest) (*TokenResponse, error) {
	user, err := s.userRepository.FindUserByEmail(req.Email)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	accessToken, accessExpiry, err := s.generateToken(user.ID, user.Email, utils.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.generateToken(user.ID, user.Email, utils.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) RefreshToken(refreshTokenStr string) (*TokenResponse, error) {
	userInfo, err := utils.ValidateToken(s.config, refreshTokenStr, utils.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Get token from DB to ensure it's still valid (hasn't been revoked)
	_, err = s.tokenRepo.GetTokenByValue(refreshTokenStr)
	if err != nil {
		return nil, utils.ErrInvalidToken
	}

	accessToken, accessExpiry, err := s.generateToken(userInfo.ID, userInfo.Email, utils.AccessToken)
	if err != nil {
		return nil, err
	}

	// Return new tokens
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr, // Keep using the same refresh token
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) Logout(userID uuid.UUID) error {
	return s.tokenRepo.DeleteUserTokens(userID, utils.RefreshToken)
}

func (s *authService) generateToken(userID uuid.UUID, email string, tokenType utils.TokenType) (string, time.Time, error) {
	var expiration time.Duration

	if tokenType == utils.AccessToken {
		expiration = s.config.JWT.AccessTokenExp
	} else {
		expiration = s.config.JWT.RefreshTokenExp
	}

	expirationTime := time.Now().Add(expiration)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"type":    string(tokenType),
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	// Store the token in the database if it's a refresh token
	if tokenType == utils.RefreshToken {
		tokenObj := &Token{
			UserID:    userID,
			Token:     tokenString,
			Type:      tokenType,
			ExpiresAt: expirationTime,
		}
		if err := s.tokenRepo.CreateToken(tokenObj); err != nil {
			return "", time.Time{}, err
		}
	}

	return tokenString, expirationTime, nil
}
