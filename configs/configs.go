package configs

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
)

//go:embed configs.yaml
var configFile string

var ConfigXx Config

func InitConfig() {
	fmt.Println(configFile)
	err := yaml.Unmarshal([]byte(configFile), &ConfigXx)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Notes  []NoteGroup  `yaml:"notes"`
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
	NoteName string `yaml:"note_name"`
	NoteKey  string `yaml:"note_key"`
	Dir      string `yaml:"dir"`
}
