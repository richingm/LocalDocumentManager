package repo

import (
	"errors"
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

func (r *ArticleRepo) Create(ctx context.Context, po *entity.ArticlePo) (*entity.ArticlePo, error) {
	err := r.db.Model(&entity.ArticlePo{}).Create(po).Error
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (r *ArticleRepo) Update(ctx context.Context, id int, fields map[string]interface{}) error {
	err := r.db.Model(&entity.ArticlePo{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleRepo) GetSort(ctx context.Context, cid int, pid int) (int64, error) {
	type countStruct struct {
		Count *int64 `gorm:"column:count"`
	}
	var count countStruct
	err := r.db.Model(&entity.ArticlePo{}).Select("max(sort) as count").Where("category_id = ? and pid = ?", cid, pid).First(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	if count.Count == nil {
		return 0, err
	}
	return *count.Count, nil
}

func (r *ArticleRepo) Delete(ctx context.Context, id int) error {
	err := r.db.Where("id = ?", id).Delete(&entity.ArticlePo{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleRepo) List(ctx context.Context, search entity.ArticleParam) ([]*entity.ArticlePo, error) {
	var res []*entity.ArticlePo
	scopes := getScopes(search)
	tx := r.db.Scopes(scopes...)
	err := tx.Order("sort asc").Find(&res).Error
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
