package belajargolanggorm

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	UserId      int    `gorm:"column:user_id"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}
