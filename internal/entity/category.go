package entity

import "time"

type CategoryPo struct {
	ID        int       `gorm:"primarykey" uri:"id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Pid       int       `gorm:"column:pid"`
	Name      string    `gorm:"column:name"`
}

func (CategoryPo) TableName() string {
	return "category"
}
