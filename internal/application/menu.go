package application

import (
	"fmt"
	"path/filepath"
	"richingm/LocalDocumentManager/internal/domain"
	"strings"
)

type MenuService struct {
}

type MenuDto struct {
	MenuName string
	MenuKey  string
	DirName  string
	DirPath  string
	Children []MenuDto
}

func NewMenuService() *MenuService {
	return &MenuService{}
}

func convertMenu(menu domain.MenuDo, rootPath string) MenuDto {
	dto := MenuDto{
		MenuName: menu.MenuName,
		DirName:  menu.DirName,
		MenuKey:  menu.MenuKey,
		DirPath:  fmt.Sprintf("%s/%s", strings.TrimRight(rootPath, "/"), menu.DirName),
	}
	for _, child := range menu.Children {
		dto.Children = append(dto.Children, convertMenu(child, rootPath))
	}
	return dto
}

func (m *MenuService) GetMenus(path string) ([]MenuDto, error) {
	rootPath := filepath.Dir(path)
	menuBiz := domain.NewMenuBiz()
	menuDos, err := menuBiz.GetMenus(path)
	if err != nil {
		return nil, err
	}
	menuDtos := make([]MenuDto, 0)
	for _, menu := range menuDos {
		menuDtos = append(menuDtos, convertMenu(menu, rootPath))
	}
	return menuDtos, err
}

func generateHTML(menuItems []MenuDto, i int) string {
	html := "<ul>"
	for _, item := range menuItems {
		if len(item.Children) > 0 {
			html += "<li class=\"lsm-sidebar-item\">"
		} else {
			html += "<li>"
		}
		icon := ""
		if i == 0 {
			icon = "icon_1"
		}
		title := fmt.Sprintf("<span>%s</span>", item.MenuName)
		if len(item.Children) > 0 {
			title = fmt.Sprintf("<i class=\"my-icon lsm-sidebar-icon %s\"></i>%s<i class=\"my-icon lsm-sidebar-more\"></i>", icon, title)
		}
		html += fmt.Sprintf("<a href=\"javascript:;\" note-key=\"%s\" title=\"%s\">%s</a>", item.MenuKey, item.MenuName, title)
		if len(item.Children) > 0 {
			html += generateHTML(item.Children, i+1)
		}
		html += "</li>"
	}
	html += "</ul>"
	return html
}

func (m *MenuService) ConvertMenusToString(list []MenuDto) string {
	return generateHTML(list, 0)
}
