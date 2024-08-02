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
	Name     string
	Children []*CategoryDto
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
	return convertCategoryDtoToString(dtos), nil
}

func generateCategoryHTML(categoryDtos []*CategoryDto, i int) string {
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

func convertCategoryDtoToString(list []*CategoryDto) string {
	return generateCategoryHTML(list, 0)
}

func convertToCategoryDto(pos []*domain.CategoryDo) []*CategoryDto {
	// Create a map to store the categories by ID
	categoryMap := make(map[int]*CategoryDto)

	// Iterate through the CategoryPo slice and create/update the CategoryDto slice
	for _, po := range pos {
		// Get or create the CategoryDto for the current CategoryPo
		dto, ok := categoryMap[po.ID]
		if !ok {
			dto = &CategoryDto{
				Id:       po.ID,
				Name:     po.Name,
				Children: make([]*CategoryDto, 0),
			}
			categoryMap[po.ID] = dto
		}

		// If the current CategoryPo has a parent, add it to the parent's children
		if po.Pid != 0 {
			parent, ok := categoryMap[po.Pid]
			if ok {
				parent.Children = append(parent.Children, dto)
			} else {
				// Create a new parent CategoryDto and add the current child to it
				parent = &CategoryDto{
					Id:       po.Pid,
					Name:     po.Name, // You may need to fetch the name of the parent category
					Children: []*CategoryDto{dto},
				}
				categoryMap[po.Pid] = parent
			}
		}
	}

	// Convert the map values to a slice
	b := make([]*CategoryDto, 0, len(categoryMap))
	for _, dto := range categoryMap {
		b = append(b, dto)
	}

	return b
}
