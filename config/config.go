package config

import (
	"github.com/randyhg/test-log-scanner/util/mylog"
	"github.com/spf13/viper"
)

type Config struct {
	ShowSql           bool        `yaml:"ShowSql"`
	MySqlUrl          string      `yaml:"MySqlUrl"`
	MySqlMaxIdle      int         `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen      int         `yaml:"MySqlMaxOpen"`
	SlaveMySqlUrl     string      `yaml:"SlaveMySqlUrl"`
	SlaveMySqlMaxIdle int         `yaml:"SlaveMySqlMaxIdle"`
	SlaveMySqlMaxOpen int         `yaml:"SlaveMySqlMaxOpen"`
	RedisCache        RedisConfig `yaml:"RedisCache"`
}

type RedisConfig struct {
	Host      []string `yaml:"Host"`      // 连接地址
	Password  string   `yaml:"Password"`  // 密码
	DB        int      `yaml:"DB"`        // 库索引
	MaxIdle   int      `yaml:"MaxIdle"`   // 最大空闲数
	MaxActive int      `yaml:"MaxActive"` // 最大连接数
}

var Instance Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		mylog.Fatalf("error reading config file, %s", err)
	}
	err = viper.Unmarshal(&Instance)
	if err != nil {
		mylog.Fatalf("unable to decode into struct, %v", err)
	}
}
