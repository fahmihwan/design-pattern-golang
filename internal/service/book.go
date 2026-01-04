package service

import (
	"best-pattern/internal/model"
	"best-pattern/internal/repository"
	"context"
	"fmt"
)

var _ BookServiceInteface = &BookService{}

type BookServiceInteface interface {
	ListBook(ctx context.Context, filters map[string]string, search string, page, limit int, sortBy, orderBy string) ([]*model.Book, int, error)
}

type BookService struct {
	repo repository.Repository
}

func NewBookService(repo repository.Repository) *BookService {
	return &BookService{
		repo: repo,
	}
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
