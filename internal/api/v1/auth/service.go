package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"

	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type AuthService interface {
	Login(c *fiber.Ctx, req LoginRequest) (*TokenResponse, error)
	RefreshToken(c *fiber.Ctx, refreshToken string) (*TokenResponse, error)
	Logout(c *fiber.Ctx, userID uuid.UUID) error
	LogoutWithToken(c *fiber.Ctx, refreshToken string) error
	Register(c *fiber.Ctx, req RegisterRequest) error
}

type authService struct {
	config         *config.Config
	userRepository user.UserRepository
	repo           AuthRepository
}

func NewAuthService(config *config.Config, userRepository user.UserRepository, repo AuthRepository) AuthService {
	return &authService{config: config, userRepository: userRepository, repo: repo}
}

func (s *authService) Login(c *fiber.Ctx, req LoginRequest) (*TokenResponse, error) {
	user, err := s.userRepository.FindUserByEmail(c, req.Email)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	accessToken, accessExpiry, err := s.generateToken(c, user.ID, user.Email, string(user.Role), utils.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.generateToken(c, user.ID, user.Email, string(user.Role), utils.RefreshToken)
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

func (s *authService) RefreshToken(c *fiber.Ctx, refreshTokenStr string) (*TokenResponse, error) {
	userInfo, err := utils.ValidateToken(s.config, refreshTokenStr, utils.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Get token from DB to ensure it's still valid (hasn't been revoked)
	_, err = s.repo.GetTokenByValue(c, refreshTokenStr)
	if err != nil {
		return nil, utils.ErrInvalidToken
	}

	accessToken, accessExpiry, err := s.generateToken(c, userInfo.ID, userInfo.Email, userInfo.Role, utils.AccessToken)
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

func (s *authService) Logout(c *fiber.Ctx, userID uuid.UUID) error {
	return s.repo.DeleteUserTokens(c, userID, utils.RefreshToken)
}

func (s *authService) LogoutWithToken(c *fiber.Ctx, refreshToken string) error {
	_, err := utils.ValidateToken(s.config, refreshToken, utils.RefreshToken)
	if err != nil {
		return err
	}

	token, err := s.repo.GetTokenByValue(c, refreshToken)
	if err != nil {
		return utils.ErrInvalidToken
	}

	return s.repo.DeleteToken(c, token.ID)
}

func (s *authService) generateToken(c *fiber.Ctx, userID uuid.UUID, email string, role string, tokenType utils.TokenType) (string, time.Time, error) {
	var expiration time.Duration

	if tokenType == utils.AccessToken {
		expiration = s.config.JWT.AccessTokenExp
	} else {
		expiration = s.config.JWT.RefreshTokenExp
	}

	expirationTime := time.Now().UTC().Add(expiration)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
		"type":    string(tokenType),
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().UTC().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	// Store the token in the redis if it's a refresh token
	if tokenType == utils.RefreshToken {
		tokenObj := &Token{
			UserID:    userID,
			Token:     tokenString,
			Type:      tokenType,
			ExpiresAt: expirationTime,
		}
		if err := s.repo.CreateToken(c, tokenObj); err != nil {
			return "", time.Time{}, err
		}
	}

	return tokenString, expirationTime, nil
}

func (s *authService) Register(c *fiber.Ctx, req RegisterRequest) error {
	newUser := &user.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      user.RoleUser,
	}
	if err := newUser.HashPassword(req.Password); err != nil {
		return err
	}

	return s.repo.CreateUser(newUser)
}
