package repo

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"richingm/LocalDocumentManager/internal/entity"
)

type ArticleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{
		db: db,
	}
}

func (r *ArticleRepo) Create(ctx context.Context, po *entity.ArticlePo) error {
	return r.db.Model(&entity.ArticlePo{}).Create(po).Error
}

func (r *ArticleRepo) List(ctx context.Context, search entity.ArticleParam) ([]*entity.ArticlePo, error) {
	var res []*entity.ArticlePo
	scopes := getScopes(search)
	tx := r.db.Scopes(scopes...)
	err := tx.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ArticleRepo) GetByID(ctx context.Context, id int) (*entity.ArticlePo, error) {
	var po entity.ArticlePo
	result := r.db.Model(&entity.ArticlePo{}).First(&po, id)
	return &po, result.Error
}

func getScopes(search entity.ArticleParam) []func(*gorm.DB) *gorm.DB {
	res := make([]func(*gorm.DB) *gorm.DB, 0)
	if search.Pid > 0 {
		res = append(res, wherePid(search.Pid))
	}
	if search.CategoryID > 0 {
		res = append(res, whereCategoryId(search.CategoryID))
	}
	return res
}

func wherePid(val int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pid = ?", val)
	}
}

func whereCategoryId(val int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category_id = ?", val)
	}
}
