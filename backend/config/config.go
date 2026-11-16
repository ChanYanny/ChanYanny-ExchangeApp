package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

var ApplicationConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	ApplicationConfig = &Config{} // 初始化配置结构体，防止空指针异常

	if err := viper.Unmarshal(ApplicationConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

}
