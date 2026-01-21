package service

import (
	"best-pattern/internal/model"
	"best-pattern/internal/repository"
	"context"
	"fmt"
	"log"
	"time"
)

var _ BookServiceInteface = &BookService{}

type BookServiceInteface interface {
	ListBook(ctx context.Context, filters map[string]string, search string, page, limit int, sortBy, orderBy string) ([]*model.Book, int, error)
	CreateBook(ctx context.Context, book *model.Book) (*model.Book, error)
	GetBookByID(ctx context.Context, id string) (*model.Book, error)
	UpdateBook(ctx context.Context, book *model.Book) (*model.Book, error)
}

type BookService struct {
	repo repository.Repository
}

func NewBookService(repo repository.Repository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) CreateBook(ctx context.Context, book *model.Book) (*model.Book, error) {

	// Generate a new UUID for the form
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	log.Print(book)
	err := s.repo.Book.Create(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("failed to create form: %w", err)
	}

	return book, nil
}

func (s *BookService) ListBook(ctx context.Context, filters map[string]string, search string, page, limit int, sortBy, orderBy string) ([]*model.Book, int, error) {

	offset := (page - 1) * limit

	books, total, err := s.repo.Book.List(ctx, repository.FilterBook{
		Pagination: repository.Pagination{
			Page:    page,
			Limit:   limit,
			Offset:  offset,
			SortBy:  sortBy,
			OrderBy: orderBy,
			Search:  search,
		},
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list forms: %w", err)
	}

	return books, total, nil

}

func (s *BookService) GetBookByID(ctx context.Context, id string) (*model.Book, error) {

	book, err := s.repo.Book.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by ID: %w", err)
	}
	return book, nil
}

func (s *BookService) UpdateBook(ctx context.Context, book *model.Book) (*model.Book, error) {
	book.UpdatedAt = time.Now()
	err := s.repo.Book.Update(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("failed to update form: %w", err)
	}
	return book, nil
}
