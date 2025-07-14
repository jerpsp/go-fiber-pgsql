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
			Email:     "john.doe@email.com",
			Password:  string(hashedPassword),
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Email:     "jane.doe@email.com",
			Password:  string(hashedPassword),
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	pool.DB.Save(users)
	fmt.Println("Users created successfully")

	return nil
}
