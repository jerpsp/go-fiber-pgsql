package development

import (
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

func Run(cfg *config.Config, pool *database.GormDB) {
	CreateUser(cfg, *pool)
	CreateBook(cfg, *pool)
}
