package development

import (
	"fmt"

	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(cfg *config.Config, pool database.GormDB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	users := []user.User{
		{
			Email:     "admin@email.com",
			Password:  string(hashedPassword),
			FirstName: "admin",
			LastName:  "admin",
			Role:      user.RoleAdmin,
			Active:    true,
		},
		{
			Email:     "moderator@email.com",
			Password:  string(hashedPassword),
			FirstName: "moderator",
			LastName:  "moderator",
			Role:      user.RoleModerator,
			Active:    true,
		},
		{
			Email:     "john.doe@email.com",
			Password:  string(hashedPassword),
			FirstName: "John",
			LastName:  "Doe",
			Role:      user.RoleUser,
			Active:    true,
		},
		{
			Email:     "jane.doe@email.com",
			Password:  string(hashedPassword),
			FirstName: "Jane",
			LastName:  "Doe",
			Role:      user.RoleUser,
			Active:    true,
		},
	}

	pool.DB.Save(users)
	fmt.Println("Users created successfully")

	return nil
}
