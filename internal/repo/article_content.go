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

func (r *ArticleContentRepo) Create(ctx context.Context, po *entity.ArticleContentPo) error {
	return r.db.Model(&entity.ArticleContentPo{}).Create(po).Error
}

func (r *ArticleContentRepo) Update(ctx context.Context, articleId int, fields map[string]interface{}) error {
	err := r.db.Model(&entity.ArticleContentPo{}).Where("article_id = ?", articleId).Updates(fields).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleContentRepo) Delete(ctx context.Context, articleId int) error {
	err := r.db.Where("article_id = ?", articleId).Delete(&entity.ArticleContentPo{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleContentRepo) GetByID(ctx context.Context, id int) (*entity.ArticleContentPo, error) {
	var po entity.ArticleContentPo
	result := r.db.Model(&entity.ArticleContentPo{}).Where("article_id = ?", id).First(&po, id)
	return &po, result.Error
}
