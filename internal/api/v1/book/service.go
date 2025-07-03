package book

import (
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
)

// Setup
type BookService interface {
	GetBooks(c *fiber.Ctx) ([]Book, error)
	GetBook(c *fiber.Ctx, bookID uuid.UUID) (Book, error)
	CreateBook(c *fiber.Ctx, newBook Book) (Book, error)
	UpdateBook(c *fiber.Ctx, bookID uuid.UUID, updatedBook BookRequest) (Book, error)
	DeleteBook(c *fiber.Ctx, bookID uuid.UUID) error
	UploadFile(c *fiber.Ctx, file *multipart.FileHeader) (string, error)
}

type bookService struct {
	config *config.Config
	repo   BookRepository
}

func NewBookService(config *config.Config, repo BookRepository) BookService {
	return &bookService{config: config, repo: repo}
}

// Service methods
func (s *bookService) GetBooks(c *fiber.Ctx) ([]Book, error) {
	books, err := s.repo.FindAllBooks(c)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *bookService) GetBook(c *fiber.Ctx, bookID uuid.UUID) (Book, error) {
	book, err := s.repo.FindBookByID(c, bookID)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (s *bookService) CreateBook(c *fiber.Ctx, newBook Book) (Book, error) {
	newBook, err := s.repo.CreateBook(c, newBook)
	if err != nil {
		return Book{}, err
	}

	return newBook, nil
}

func (s *bookService) UpdateBook(c *fiber.Ctx, bookID uuid.UUID, updateBook BookRequest) (Book, error) {
	book, err := s.repo.FindBookByID(c, bookID)
	if err != nil {
		return Book{}, err
	}
	book.Title = updateBook.Title
	book.Author = updateBook.Author
	book.UpdatedAt = time.Now()

	updatedBook, err := s.repo.UpdateBook(c, book)
	if err != nil {
		return Book{}, err
	}

	return updatedBook, nil
}

func (s *bookService) DeleteBook(c *fiber.Ctx, bookID uuid.UUID) error {
	return s.repo.DeleteBook(c, bookID)
}

func (s *bookService) UploadFile(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	if err := c.SaveFile(file, "./public/uploads/"+file.Filename); err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to save file")
	}

	return file.Filename, nil
}
