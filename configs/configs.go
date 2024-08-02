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
	Server      ServerConfig `yaml:"server"`
	MysqlConfig MysqlConfig  `yaml:"mysql"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

type MysqlConfig struct {
	Path            string `yaml:"path"`              // 服务器地址:端口
	Config          string `yaml:"config"`            // 高级配置
	Dbname          string `yaml:"db-name"`           // 数据库名
	Username        string `yaml:"username"`          // 数据库用户名
	Password        string `yaml:"password"`          // 数据库密码
	MaxIdleConns    int    `yaml:"max-idle-conns"`    // 空闲中的最大连接数
	MaxOpenConns    int    `yaml:"max-open-conns"`    // 打开到数据库的最大连接数
	ConnMaxLifetime int    `yaml:"conn-max-lifetime"` // 连接最大生存时间, 单位: 分钟
}
