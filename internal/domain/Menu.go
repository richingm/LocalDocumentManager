package domain

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type MenuDo struct {
	MenuName string   `yaml:"menu_name,omitempty"`
	DirName  string   `yaml:"dir_name,omitempty"`
	MenuKey  string   `yaml:"menu_key,omitempty"`
	Children []MenuDo `yaml:"children,omitempty"`
}

type MenuBiz struct {
}

func NewMenuBiz() *MenuBiz {
	return &MenuBiz{}
}

func (m *MenuBiz) GetMenus(path string) ([]MenuDo, error) {
	// 读取YAML文件内容
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 解析YAML文件
	var menuDo []MenuDo
	err = yaml.Unmarshal(yamlFile, &menuDo)
	if err != nil {
		return nil, err
	}

	return menuDo, nil
}
