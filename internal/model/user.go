package model

type User struct {
	ID    string `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(255);not null"`
	Email string `gorm:"type:varchar(255);not null"`
}

func (User) TableName() string {
	return "users"
}
