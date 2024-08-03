package entity

type ArticleContentPo struct {
	ID        int     `gorm:"primarykey" uri:"id"`
	CreatedAt []uint8 `gorm:"column:created_at"`
	ArticleID int     `gorm:"column:article_id"`
	Content   string  `gorm:"column:content;type:mediumtext"`
}

func (ArticleContentPo) TableName() string {
	return "article_content"
}
