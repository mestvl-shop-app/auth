package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	GrpcServer GrpcServer `yaml:"grpc_server" env-required:"true"`
	Database   Database   `yaml:"database" env-required:"true"`
}

type GrpcServer struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type Database struct {
	Net                string        `yaml:"net" env-default:"tcp"`
	Host               string        `yaml:"host" env-required:"true"`
	Port               string        `yaml:"port" env-required:"true"`
	DBName             string        `yaml:"db_name" env-required:"true"`
	User               string        `yaml:"user" env-required:"true"`
	Password           string        `yaml:"password" env-required:"true"`
	SSLMode            string        `yaml:"sslmode" env-required:"true"`
	TimeZone           string        `yaml:"time_zone"`
	Timeout            time.Duration `yaml:"timeout" env-default:"2s"`
	MaxIdleConnections int           `yaml:"max_idle_connections" env-default:"40"`
	MaxOpenConnections int           `yaml:"max_open_connections" env-default:"40"`
}

func MustLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("load .env failed: %s", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	return &cfg
}
