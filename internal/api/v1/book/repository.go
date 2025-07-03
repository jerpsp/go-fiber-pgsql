package book

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

type BookRepository interface {
	FindAllBooks(c *fiber.Ctx) ([]Book, error)
	FindBookByID(c *fiber.Ctx, bookID uuid.UUID) (Book, error)
	CreateBook(c *fiber.Ctx, newBook Book) (Book, error)
	UpdateBook(c *fiber.Ctx, updatedBook Book) (Book, error)
	DeleteBook(c *fiber.Ctx, bookID uuid.UUID) error
}

type bookRepository struct {
	config *config.Config
	db     *database.GormDB
}

func NewBookRepository(cfg *config.Config, db *database.GormDB) BookRepository {
	return &bookRepository{config: cfg, db: db}
}

// Repository methods

func (r *bookRepository) FindAllBooks(c *fiber.Ctx) ([]Book, error) {
	var books []Book
	if err := r.db.DB.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) FindBookByID(c *fiber.Ctx, bookID uuid.UUID) (Book, error) {
	var book Book
	if err := r.db.DB.Where("id = ?", bookID).First(&book).Error; err != nil {
		return Book{}, err
	}
	return book, nil
}

func (r *bookRepository) CreateBook(c *fiber.Ctx, newBook Book) (Book, error) {
	if err := r.db.DB.Create(&newBook).Error; err != nil {
		return Book{}, err
	}
	return newBook, nil
}

func (r *bookRepository) UpdateBook(c *fiber.Ctx, updatedBook Book) (Book, error) {
	if err := r.db.DB.Save(&updatedBook).Error; err != nil {
		return Book{}, err
	}

	return updatedBook, nil
}

func (r *bookRepository) DeleteBook(c *fiber.Ctx, bookID uuid.UUID) error {
	if err := r.db.DB.Delete(&Book{}, bookID).Error; err != nil {
		return err
	}
	return nil
}
