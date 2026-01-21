package model

import (
	"time"
)

type Book struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Author      string     `gorm:"type:varchar(255);not null"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:timestamp" json:"deleted_at,omitempty"`
}

func (Book) TableName() string {
	return "books"
}
