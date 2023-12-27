package configs

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed configs.yaml
var configFile string

var ConfigXx Config

func InitConfig() {
	err := yaml.Unmarshal([]byte(configFile), &ConfigXx)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Server              ServerConfig `yaml:"server"`
	Docs                []string     `yaml:"docs"`
	Notes               []NoteGroup  `yaml:"notes"`
	DefaultDisplayLevel int64        `yaml:"default_display_level"`
	MenusYamlFile       string       `yaml:"menus_yaml_file"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

type NoteGroup struct {
	GroupName string      `yaml:"note_group"`
	Children  []NoteChild `yaml:"children"`
}

type NoteChild struct {
	NoteName     string `yaml:"note_name"`
	NoteKey      string `yaml:"note_key"`
	Dir          string `yaml:"dir"`
	DisplayLevel int64  `yaml:"display_level"`
}
