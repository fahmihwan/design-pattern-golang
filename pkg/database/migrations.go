package database

import (
	"time"

	"gorm.io/gorm"
)

// Migration represents an individual migration record
type Migration struct {
	ID    uint      `gorm:"primaryKey"`
	Name  string    `gorm:"unique;not null"`
	RunAt time.Time `gorm:"not null"`
}

// AutoMigrate ensures that the migrations table is present
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Migration{})
}
