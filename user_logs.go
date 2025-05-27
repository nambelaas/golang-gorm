package belajargolanggorm

type UserLog struct {
	ID        int       `gorm:"column:id;primary_key;autoIncrement`
	UserId    int       `gorm:"column:user_id"`
	Action    string    `gorm:"column:action"`
	CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli;<-:create"`
	UpdatedAt int64 `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}
