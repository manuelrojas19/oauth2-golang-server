package entities

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Username  string    `gorm:"type:varchar(255);unique;not null"`
	Email     string    `gorm:"type:varchar(255);unique"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}
