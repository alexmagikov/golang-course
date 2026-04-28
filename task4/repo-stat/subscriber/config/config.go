package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type GitHub struct {
	Token string `yaml:"token" env:"GITHUB_TOKEN" env-default:""`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"myuser"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"mypassword"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"mydb"`
}

type Config struct {
	App      App               `yaml:"app"`
	Database Database          `yaml:"database"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	Github   GitHub            `yaml:"github"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
