package domain

import (
	"golang.org/x/net/context"
	"richingm/LocalDocumentManager/internal/repo"
	"time"
)

type CategoryBiz struct {
	categoryRepo *repo.CategoryRepo
}

type CategoryDo struct {
	ID        int       `gorm:"primarykey" uri:"id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Pid       int       `gorm:"column:pid"`
	Name      string    `gorm:"column:name"`
}

func NewCategoryBiz(ctx context.Context, categoryRepo *repo.CategoryRepo) *CategoryBiz {
	return &CategoryBiz{
		categoryRepo: categoryRepo,
	}
}

func (b *CategoryBiz) Get(ctx context.Context, id int) (*CategoryDo, error) {
	po, err := b.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	do := &CategoryDo{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		Pid:       po.Pid,
		Name:      po.Name,
	}
	return do, err
}

func (b *CategoryBiz) List(ctx context.Context) ([]*CategoryDo, error) {
	list, err := b.categoryRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*CategoryDo, 0)
	for _, po := range list {
		res = append(res, &CategoryDo{
			ID:        po.ID,
			CreatedAt: po.CreatedAt,
			Pid:       po.Pid,
			Name:      po.Name,
		})
	}
	return res, nil
}
