package repo

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"richingm/LocalDocumentManager/internal/entity"
)

type ArticleContentRepo struct {
	db *gorm.DB
}

func NewArticleContentRepo(db *gorm.DB) *ArticleContentRepo {
	return &ArticleContentRepo{
		db: db,
	}
}

func (r *ArticleContentRepo) Create(ctx context.Context, po *entity.ArticlePo) error {
	return r.db.Model(&entity.ArticleContentPo{}).Create(po).Error
}

func (r *ArticleContentRepo) GetByID(ctx context.Context, id int) (*entity.ArticleContentPo, error) {
	var po entity.ArticleContentPo
	result := r.db.Model(&entity.ArticleContentPo{}).Where("article_id = ?", id).First(&po, id)
	return &po, result.Error
}
