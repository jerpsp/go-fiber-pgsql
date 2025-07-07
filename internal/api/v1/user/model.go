package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid; default:uuid_generate_v4()" json:"id"`
	Email     string    `gorm:"size:255;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255" json:"password"`
	FirstName string    `gorm:"size:100" json:"first_name"`
	LastName  string    `gorm:"size:100" json:"last_name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
