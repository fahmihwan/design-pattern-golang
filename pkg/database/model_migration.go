package database

import "time"

type Book struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Author      string     `gorm:"type:varchar(255);not null" json:"author"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:timestamptz" json:"deleted_at,omitempty"`
}
