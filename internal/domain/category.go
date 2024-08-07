package domain

import (
	"golang.org/x/net/context"
	"richingm/LocalDocumentManager/internal/entity"
	"richingm/LocalDocumentManager/internal/repo"
)

type CategoryBiz struct {
	categoryRepo *repo.CategoryRepo
}

type CategoryDo struct {
	ID        int
	CreatedAt []uint8
	Pid       int
	Name      string
	Children  []CategoryDo
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

func (b *CategoryBiz) GetByPid(ctx context.Context, pid int) ([]CategoryDo, error) {
	list, err := b.categoryRepo.GetByPid(ctx, pid)
	if err != nil {
		return nil, err
	}
	res := make([]CategoryDo, 0)
	for _, po := range list {
		res = append(res, CategoryDo{
			ID:        po.ID,
			CreatedAt: po.CreatedAt,
			Pid:       po.Pid,
			Name:      po.Name,
		})
	}
	return res, err
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

func (b *CategoryBiz) Create(ctx context.Context, pid int, title string, content string) error {
	sortNum, err := b.categoryRepo.GetSort(ctx, pid)
	if err != nil {
		return err
	}

	err = b.categoryRepo.Create(ctx, &entity.CategoryPo{
		Pid:  pid,
		Name: title,
		Sort: int(sortNum + 1),
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *CategoryBiz) Delete(ctx context.Context, id int) error {
	err := b.categoryRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
