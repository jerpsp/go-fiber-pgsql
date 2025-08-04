package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type Token struct {
	ID        uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	Token     string          `gorm:"size:1024;not null" json:"token"`
	Type      utils.TokenType `gorm:"size:20;not null" json:"type"`
	ExpiresAt time.Time       `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
