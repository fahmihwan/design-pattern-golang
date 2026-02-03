package repository

import "gorm.io/gorm"

type Repository struct {
	Book BookRepo
	User UserRepo
}

type Pagination struct {
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	SortBy  string `json:"sort_by,omitempty"`
	OrderBy string `json:"order_by,omitempty"`
	Search  string `json:"search,omitempty"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Book: NewBookRepository(db),
		User: NewUserRepository(db),
	}
}
