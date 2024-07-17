package entities

import "time"

type BaseGormEntity struct {
	Id        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
