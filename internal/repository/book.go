package repository

import (
	"best-pattern/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookRepository struct {
	db *gorm.DB
}

type BookRepo interface {
	Create(ctx context.Context, book *model.Book) error
	List(ctx context.Context, filter FilterBook) (res []*model.Book, total int, err error)
	setFilter(db *gorm.DB, filter FilterBook) *gorm.DB
}

var _ BookRepo = (*BookRepository)(nil)

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

type FilterBook struct {
	Pagination
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Author      string     `json:"author,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func (r *BookRepository) Create(ctx context.Context, book *model.Book) error {

	err := r.db.WithContext(ctx).Create(book).Error
	if err != nil {
		return fmt.Errorf("failed to create form: %w", err)
	}
	return nil
}

func (r *BookRepository) List(ctx context.Context, filter FilterBook) (res []*model.Book, total int, err error) {

	funcName := "List"
	tableName := model.Book{}.TableName()

	// Pastikan slice kosong (bukan nil)
	res = make([]*model.Book, 0)

	// GORM pakai context (ini bukan OpenTelemetry; ini untuk cancel/timeout dari request)
	db := r.db.WithContext(ctx)

	// Operation `count`
	var count int64
	err = r.setFilter(db, filter).Model(&model.Book{}).Count(&count).Error

	if err != nil {
		return res, total, fmt.Errorf("failed to %s %s count: %w", funcName, tableName, err)
	}

	if count == 0 {
		return
	}
	total = int(count)

	if filter.SortBy == "" {
		filter.SortBy = "id"
	}

	order := strings.ToUpper(filter.OrderBy)
	desc := order == "DESC" // default ASC kalau kosong / selain DESC

	// Operation `select`
	if err = r.setFilter(db, filter).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: filter.SortBy},
			Desc:   desc,
		}).
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&res).Error; err != nil {
		return res, total, fmt.Errorf("failed to %s %s find: %w", funcName, tableName, err)
	}

	return res, total, nil
}

func (r *BookRepository) setFilter(db *gorm.DB, filter FilterBook) *gorm.DB {
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		db = db.Where("title ILIKE ? OR author ILIKE ?", like, like)

	}

	if filter.ID != "" {
		db = db.Where("id = ?", filter.ID)
	}

	// Soft delete: hanya yang belum dihapus
	db = db.Where("deleted_at IS NULL")

	return db
}
