package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer `yaml:"http_server"`
	DB         `yaml:"db"`
	Env        string `yaml:"env" env-default:"local"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type HttpServer struct {
	Adress      string        `yaml:"adress" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func Require() *Config {
	// confPath := os.Args[1]
	confPath := "/Users/amir/GoProjects/URLShortener/config/local.yaml"
	if confPath == "" {
		log.Fatal("no config file. Pass the config file path as run parameter")
	}
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		log.Fatalf("file %s does not exist", confPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(confPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}
	return &cfg
}
