package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole represents user roles in the system
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid; default:uuid_generate_v4()" json:"id"`
	Email        string    `gorm:"size:255;uniqueIndex" json:"email"`
	Password     string    `gorm:"size:255" json:"-"`
	FirstName    string    `gorm:"size:100" json:"first_name"`
	LastName     string    `gorm:"size:100" json:"last_name"`
	Role         UserRole  `gorm:"size:20;default:'user'" json:"role"`
	ProfileImage string    `gorm:"size:255" json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HasRole checks if user has the specified role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// HasAnyRole checks if user has any of the specified roles
func (u *User) HasAnyRole(roles ...UserRole) bool {
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}
