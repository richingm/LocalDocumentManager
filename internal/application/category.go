package application

import (
	"fmt"
	"golang.org/x/net/context"
	"richingm/LocalDocumentManager/internal/domain"
	"richingm/LocalDocumentManager/internal/infrastructure/mysql"
	"richingm/LocalDocumentManager/internal/repo"
)

type CategoryService struct {
}

type CategoryDto struct {
	Id       int
	Pid      int
	Name     string
	Sort     int
	Children []CategoryDto
}

func NewCategoryService(ctx context.Context) *CategoryService {
	return &CategoryService{}
}

func (r *CategoryService) Create(ctx context.Context, pid int, title string, content string) error {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	err := categoryBiz.Create(ctx, pid, title, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) Update(ctx context.Context, id int, title string, orderSort int) error {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	err := categoryBiz.Update(ctx, id, title, orderSort)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	err := categoryBiz.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) Get(ctx context.Context, id int) (*CategoryDto, error) {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	do, err := categoryBiz.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &CategoryDto{
		Id:   do.ID,
		Name: do.Name,
		Sort: do.Sort,
	}, nil
}

func (r *CategoryService) ListHtml(ctx context.Context) (string, error) {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	list, err := categoryBiz.List(ctx)
	if err != nil {
		return "", err
	}
	dtos := convertToCategoryDto(list)
	dtos = buildTree(dtos, 0)
	return convertCategoryDtoToString(dtos), nil
}

func (s *CategoryService) Nodes(ctx context.Context, categoryId int) (NodeDto, error) {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	categoryDo, err := categoryBiz.Get(ctx, categoryId)
	if err != nil {
		return NodeDto{}, err
	}
	res := NodeDto{
		ID:       fmt.Sprintf("%d", categoryDo.ID),
		Topic:    categoryDo.Name,
		Expanded: false,
	}
	categoryDos, err := categoryBiz.GetByPid(ctx, categoryId)
	if err != nil {
		return NodeDto{}, err
	}

	res.Children = buildTreeCategory(categoryDos)
	return res, nil
}

func buildTreeCategory(categoryDos []domain.CategoryDo) []*NodeDto {
	res := make([]*NodeDto, 0, len(categoryDos))
	for _, categoryDo := range categoryDos {
		res = append(res, buildTreeCategoryLoop(categoryDo))
	}
	return res
}

func buildTreeCategoryLoop(do domain.CategoryDo) *NodeDto {
	if len(do.Children) == 0 {
		return &NodeDto{
			ID:       fmt.Sprintf("%d", do.ID),
			Topic:    do.Name,
			Children: nil,
			Expanded: false,
		}
	}
	children := make([]*NodeDto, 0, len(do.Children))
	for _, child := range do.Children {
		children = append(children, buildTreeCategoryLoop(child))
	}
	return &NodeDto{
		ID:       fmt.Sprintf("%d", do.ID),
		Topic:    do.Name,
		Children: children,
		Expanded: false,
	}
}

func generateCategoryHTML(categoryDtos []CategoryDto, i int) string {
	html := "<ul>"
	for _, categoryDto := range categoryDtos {
		if len(categoryDto.Children) > 0 {
			html += "<li class=\"lsm-sidebar-item\">"
		} else {
			html += "<li>"
		}
		icon := ""
		if i == 0 {
			icon = "icon_1"
		}
		title := fmt.Sprintf("<span>%s</span>", categoryDto.Name)
		if len(categoryDto.Children) > 0 {
			title = fmt.Sprintf("<i class=\"my-icon lsm-sidebar-icon %s\"></i>%s<i class=\"my-icon lsm-sidebar-more\"></i>", icon, title)
		}
		html += fmt.Sprintf("<a href=\"javascript:;\" note-id=\"%d\" title=\"%s\">%s</a>", categoryDto.Id, categoryDto.Name, title)
		if len(categoryDto.Children) > 0 {
			html += generateCategoryHTML(categoryDto.Children, i+1)
		}
		html += "</li>"
	}
	html += "</ul>"
	return html
}

func convertCategoryDtoToString(list []CategoryDto) string {
	return generateCategoryHTML(list, 0)
}

func convertToCategoryDto(pos []*domain.CategoryDo) []CategoryDto {
	res := make([]CategoryDto, 0)
	for _, po := range pos {
		res = append(res, CategoryDto{
			Id:       po.ID,
			Pid:      po.Pid,
			Name:     po.Name,
			Children: make([]CategoryDto, 0),
		})
	}
	return res
}

func buildTree(items []CategoryDto, pid int) []CategoryDto {
	var result []CategoryDto
	for _, item := range items {
		if item.Pid == pid {
			children := buildTree(items, item.Id)
			item.Children = children
			result = append(result, item)
		}
	}
	return result
}
