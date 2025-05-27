package belajargolanggorm

import (
	"time"
)

type Wallet struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserId    int       `gorm:"column:user_id"`
	Balance   float64   `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	User      *User     `gorm:"foreignKey:user_id;references:id"` //relasi one to one
}
