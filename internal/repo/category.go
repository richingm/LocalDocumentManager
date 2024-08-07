package repo

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"richingm/LocalDocumentManager/internal/entity"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (r *CategoryRepo) Create(ctx context.Context, po *entity.CategoryPo) error {
	return r.db.Model(&entity.CategoryPo{}).Create(po).Error
}

func (r *CategoryRepo) Update(ctx context.Context, po *entity.CategoryPo) error {
	return r.db.Model(&entity.CategoryPo{}).Save(po).Error
}

func (r *CategoryRepo) Delete(ctx context.Context, id int) error {
	return r.db.Model(&entity.CategoryPo{}).Delete(&entity.CategoryPo{ID: id}).Error
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int) (*entity.CategoryPo, error) {
	var category entity.CategoryPo
	result := r.db.Model(&entity.CategoryPo{}).First(&category, id)
	return &category, result.Error
}

func (r *CategoryRepo) GetByPid(ctx context.Context, pid int) ([]entity.CategoryPo, error) {
	list := make([]entity.CategoryPo, 0)
	result := r.db.Model(&entity.CategoryPo{}).Where("pid = ?", pid).Find(&list)
	return list, result.Error
}

func (r *CategoryRepo) List(ctx context.Context) ([]*entity.CategoryPo, error) {
	var categories []*entity.CategoryPo
	result := r.db.Model(&entity.CategoryPo{}).Find(&categories).Order("sort asc")
	return categories, result.Error
}
