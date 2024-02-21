package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPC     GRPC
	Database Database
}

type GRPC struct {
	Port string
}

type Database struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	Timezone string
}

func NewConfig() *Config {
	return &Config{
		GRPC: GRPC{
			Port: viper.GetString("grpc.port"),
		},
		Database: Database{
			Host:     viper.GetString("postgres.host"),
			User:     viper.GetString("postgres.user"),
			Password: viper.GetString("postgres.password"),
			DBName:   viper.GetString("postgres.dbname"),
			Port:     viper.GetString("postgres.port"),
			SSLMode:  viper.GetString("postgres.sslmode"),
			Timezone: viper.GetString("postgres.timezone"),
		},
	}
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
