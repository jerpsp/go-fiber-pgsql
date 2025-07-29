package production

import (
	"fmt"

	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/book"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

func CreateBook(cfg *config.Config, pool database.GormDB) error {
	books := []book.Book{
		{
			Title:  "The Great Gatsby",
			Author: "F. Scott Fitzgerald",
		},
		{
			Title:  "To Kill a Mockingbird",
			Author: "Harper Lee",
		},
		{
			Title:  "1984",
			Author: "George Orwell",
		},
		{
			Title:  "Pride and Prejudice",
			Author: "Jane Austen",
		},
	}

	pool.DB.Save(books)
	fmt.Println("Books created successfully")

	return nil
}
