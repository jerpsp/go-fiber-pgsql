package auth

import (
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type AuthRepository interface {
	CreateToken(token *Token) error
	GetTokenByValue(tokenStr string) (*Token, error)
	DeleteToken(tokenID uuid.UUID) error
	DeleteUserTokens(userID uuid.UUID, tokenType utils.TokenType) error
}

type authRepository struct {
	config *config.Config
	db     *database.GormDB
	redis  *database.RedisDB
}

func NewAuthRepository(config *config.Config, db *database.GormDB, redis *database.RedisDB) AuthRepository {
	return &authRepository{config: config, db: db, redis: redis}
}

func (r *authRepository) CreateToken(token *Token) error {
	// location, _ := time.LoadLocation("Asia/Bangkok")
	// expire := time.Unix(token.ExpiresAt.Unix(), 0).In(location)
	// err := r.redis.Client.Set(context.Background(), "token", token.Token, expire.Sub(token.CreatedAt.In(location))).Err()
	// if err != nil {
	// 	return err
	// }
	return r.db.DB.Create(token).Error
}

func (r *authRepository) GetTokenByValue(tokenStr string) (*Token, error) {
	var token Token
	if err := r.db.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authRepository) DeleteToken(tokenID uuid.UUID) error {
	return r.db.DB.Where("id = ?", tokenID).Delete(&Token{}).Error
}

func (r *authRepository) DeleteUserTokens(userID uuid.UUID, tokenType utils.TokenType) error {
	return r.db.DB.Where("user_id = ? AND type = ?", userID, tokenType).Delete(&Token{}).Error
}
