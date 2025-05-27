package belajargolanggorm

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int       `gorm:"column:id;primaryKey;<-:create"`
	Name         Name      `gorm:"embedded"`
	Password     string    `gorm:"column:password"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Information  string    `gorm:"-"`
	Wallet       Wallet    `gorm:"foreignKey:user_id;references:id"` //relasi one to one
	Addresses    []Address `gorm:"foreignKey:user_id;references:id"` //relasi one to many
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_id"`
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	MiddleName string `gorm:"column:middle_name"`
}

// hook before create
func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.Name.FirstName == "" {
		u.Name.FirstName = "Test hook - " + time.Now().Format("20250511090001")
	}
	return nil
}
