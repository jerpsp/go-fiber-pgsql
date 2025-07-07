package seeds

import (
	"fmt"

	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/db/seeds/development"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

func CreateSeedData(cfg *config.Config, pool *database.GormDB) {
	switch cfg.Server.ENV {
	case "development":
		fmt.Println("Seeding development data...")
		development.Run(cfg, pool)
	}

	fmt.Println("Seeding completed.")
}
