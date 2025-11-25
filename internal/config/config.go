package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env-default:"dev" env-required:"true"`
	Server   `yaml:"server"`
	Database `yaml:"database"`
}

type Server struct {
	Port string `yaml:"port" env-default:":8080"`
}

type Database struct {
	Port       string `yaml:"port"`
	DBhost     string `yaml:"host"`
	DBname     string `yaml:"db_name"`
	DBpassword string `yaml:"db_password"`
	SSLMode    string `yaml:"sslmode"`
	DBusername string `yaml:"db_username"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		log.Fatal("CONFIG variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("No .yml file: %s", err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read .yml: %s", err.Error())
	}

	return &cfg
}
