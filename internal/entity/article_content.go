package entity

import "time"

type ArticleContentPo struct {
	ID        int       `gorm:"primarykey" uri:"id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	ArticleID int       `gorm:"column:article_id"`
	Content   string    `gorm:"column:content;type:mediumtext"`
}

func (ArticleContentPo) TableName() string {
	return "article_content"
}
