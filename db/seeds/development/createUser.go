package development

import (
	"fmt"

	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

func CreateUser(cfg *config.Config, pool database.GormDB) error {
	users := []user.User{
		{
			Email:     "john.doe@email.com",
			Password:  "password",
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Email:     "jane.doe@email.com",
			Password:  "password",
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	pool.DB.Save(users)
	fmt.Println("Users created successfully")

	return nil
}
