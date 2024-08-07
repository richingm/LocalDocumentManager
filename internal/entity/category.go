package entity

type CategoryPo struct {
	ID        int     `gorm:"primarykey" uri:"id"`
	CreatedAt []uint8 `gorm:"column:created_at"`
	Pid       int     `gorm:"column:pid"`
	Name      string  `gorm:"column:name"`
	Sort      int     `gorm:"sort"`
}

func (CategoryPo) TableName() string {
	return "category"
}
