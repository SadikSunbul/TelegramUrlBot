package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
)

type Schema struct {
	MongoDbConnect string `env:"mongoDbConnect"`
	BootIdTelegram string `env:"bootIdTelegram"`
	DbName         string `env:"dbName"`
	ApiDomain      string `env :"apiDomain"`
}

var (
	cfg Schema
)

func LoadConfig(configname string) *Schema {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	err := godotenv.Load(filepath.Join(currentDir, configname))

	if err != nil {
		log.Printf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}

func GetConfig() *Schema {
	return &cfg
}
