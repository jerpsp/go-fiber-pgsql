package book

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID `gorm:"type:uuid; default:uuid_generate_v4()" json:"id"`
	Title     string    `gorm:"size:255" json:"title"`
	Author    string    `gorm:"size:255" json:"author"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
