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
	Children []CategoryDto
}

func NewCategoryService(ctx context.Context) *CategoryService {
	return &CategoryService{}
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
