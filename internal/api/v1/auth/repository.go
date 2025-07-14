package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type AuthRepository interface {
	CreateToken(c *fiber.Ctx, token *Token) error
	GetTokenByValue(c *fiber.Ctx, tokenStr string) (*Token, error)
	DeleteToken(c *fiber.Ctx, tokenID uuid.UUID) error
	DeleteUserTokens(c *fiber.Ctx, userID uuid.UUID, tokenType utils.TokenType) error
	CreateUser(user *user.User) error
}

type authRepository struct {
	config *config.Config
	db     *database.GormDB
	redis  *database.RedisDB
}

func NewAuthRepository(config *config.Config, db *database.GormDB, redis *database.RedisDB) AuthRepository {
	return &authRepository{config: config, redis: redis, db: db}
}

func (r *authRepository) CreateToken(c *fiber.Ctx, token *Token) error {
	expire := token.ExpiresAt.UTC()
	now := time.Now().UTC()

	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}

	// Create Redis key with token ID
	tokenKey := fmt.Sprintf("token:%s", token.ID.String())

	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = r.redis.Client.Set(context.Background(), tokenKey, tokenData, expire.Sub(now)).Err()
	if err != nil {
		return err
	}

	tokenLookupKey := fmt.Sprintf("token_lookup:%s", token.Token)
	err = r.redis.Client.Set(context.Background(), tokenLookupKey, token.ID.String(), expire.Sub(now)).Err()
	if err != nil {
		return err
	}

	userTokenKey := fmt.Sprintf("user_tokens:%s:%s", token.UserID.String(), token.Type)
	err = r.redis.Client.SAdd(context.Background(), userTokenKey, token.ID.String()).Err()
	if err != nil {
		return err
	}

	r.redis.Client.ExpireAt(context.Background(), userTokenKey, expire)

	return nil
}

func (r *authRepository) GetTokenByValue(c *fiber.Ctx, tokenStr string) (*Token, error) {
	tokenLookupKey := fmt.Sprintf("token_lookup:%s", tokenStr)
	tokenID, err := r.redis.Client.Get(context.Background(), tokenLookupKey).Result()
	if err != nil {
		return nil, err
	}

	tokenKey := fmt.Sprintf("token:%s", tokenID)
	tokenData, err := r.redis.Client.Get(context.Background(), tokenKey).Bytes()
	if err != nil {
		return nil, err
	}

	var token Token
	if err := json.Unmarshal(tokenData, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *authRepository) DeleteToken(c *fiber.Ctx, tokenID uuid.UUID) error {
	tokenKey := fmt.Sprintf("token:%s", tokenID.String())
	tokenData, err := r.redis.Client.Get(context.Background(), tokenKey).Bytes()
	if err != nil {
		return err
	}

	var token Token
	if err := json.Unmarshal(tokenData, &token); err != nil {
		return err
	}

	tokenLookupKey := fmt.Sprintf("token_lookup:%s", token.Token)
	if err := r.redis.Client.Del(context.Background(), tokenLookupKey).Err(); err != nil {
		return err
	}

	userTokenKey := fmt.Sprintf("user_tokens:%s:%s", token.UserID.String(), token.Type)
	if err := r.redis.Client.SRem(context.Background(), userTokenKey, tokenID.String()).Err(); err != nil {
		return err
	}

	return r.redis.Client.Del(context.Background(), tokenKey).Err()
}

func (r *authRepository) DeleteUserTokens(c *fiber.Ctx, userID uuid.UUID, tokenType utils.TokenType) error {
	// Get all token IDs for this user and type
	userTokenKey := fmt.Sprintf("user_tokens:%s:%s", userID.String(), tokenType)
	tokenIDs, err := r.redis.Client.SMembers(context.Background(), userTokenKey).Result()
	if err != nil {
		return err
	}

	// Delete each token
	for _, tokenID := range tokenIDs {
		tokenKey := fmt.Sprintf("token:%s", tokenID)
		tokenData, err := r.redis.Client.Get(context.Background(), tokenKey).Bytes()
		if err != nil {
			continue
		}

		var token Token
		if err := json.Unmarshal(tokenData, &token); err != nil {
			continue
		}

		tokenLookupKey := fmt.Sprintf("token_lookup:%s", token.Token)

		r.redis.Client.Del(context.Background(), tokenLookupKey)
		r.redis.Client.Del(context.Background(), tokenKey)
	}

	return r.redis.Client.Del(context.Background(), userTokenKey).Err()
}

func (r *authRepository) CreateUser(user *user.User) error {
	return r.db.DB.Create(user).Error
}
