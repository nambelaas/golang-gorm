package belajargolanggorm

import "time"

type Address struct {
	ID        int       `gorm:"column:id;primary_key"`
	UserId    int       `gorm:"column:user_id"`
	Address   string    `gorm:"column:address"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      User      `gorm:"foreignKey:user_id;references:id"` //relasi many to one
}

func (a *Address) TableName() string {
	return "addresses"
}
