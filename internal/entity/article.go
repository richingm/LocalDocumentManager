package entity

import "time"

type ArticlePo struct {
	ID         int       `gorm:"primarykey" uri:"id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	Pid        int       `gorm:"column:pid"`
	CategoryID int       `gorm:"column:category_id"`
	Title      string    `gorm:"column:title"`
}

func (ArticlePo) TableName() string {
	return "article"
}

type ArticleParam struct {
	CategoryID int `gorm:"column:category_id"`
	Pid        int `gorm:"column:pid"`
}
