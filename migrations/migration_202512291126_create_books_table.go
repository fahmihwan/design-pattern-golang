package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migration_202512291126() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202512291126",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS books (
                    id BIGSERIAL PRIMARY KEY,
					title VARCHAR(255) NOT NULL,
					author VARCHAR(255) NOT NULL,
					description TEXT,
					created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
					deleted_at TIMESTAMPTZ
                );

				CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON books (deleted_at);
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`
				DROP TABLE IF EXISTS books
			`).Error
		},
	}
}
