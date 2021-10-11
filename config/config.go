package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		USERNAME string `env:"DB_USERNAME"`
		PASSWORD string `env:"DB_PASSWORD"`
		NETWORK  string `env:"DB_NETWORK" envDefault:"tcp"`
		SERVER   string `env:"DB_ADDRESS"`
		PORT     string `env:"DB_PORT" envDefault:"3306"`
		DATABASE string `env:"DB_DATABASE"`
	}
	Server struct {
		PORT string `env:"PORT"`
	}
	HashID struct {
		Salt      string `env:"HASHID_SALT"`
		MinLength int    `env:"HASHID_MIN_LENGTH"`
	}
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}
	// log.Printf("%+v\n", config)
	return config
}
