package config

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     string `envconfig:"DB_PORT" required:"true"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBName     string `envconfig:"DB_NAME" required:"true"`
	DBSSLMode  string `envconfig:"DB_SSL_MODE" required:"true"`
}

func GetConfig() *Config {
	// initialize config struct and parse environment variables
	once.Do(func() {
		// make config struct singleton
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Configuration loaded successfully.")
	})

	return &config
}
