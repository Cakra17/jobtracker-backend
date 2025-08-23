package config

import (
	"log"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Name string `env:"DB_NAME"`
	Port string `env:"DB_PORT"`
	Host string `env:"DB_HOST"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	MaxOpenConnection int `env:"DB_MAX_OPEN_CONNECTION,default=0"`
	MaxIdleConnection int `env:"DB_MAX_IDLE_CONNECTION,default=0"`
	MaxConnLifetime int `env:"DB_MAX_CONN_LIFETIME,default=30"`
	MaxConnIdleTime int `env:"DB_MAX_CONN_IDLE_TIME,default=1"`
}

type Config struct {
	Database DatabaseConfig
	AppPort string `env:"APP_PORT,default=6969"`
	Env string `env:"ENV"`

	JWTSecret string `env:"JWT_SECRET"`
	BcryptSalt int `env:"BCRYPT_SALT"`
}

func InitializeConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	var cfg Config
	err = envdecode.Decode(&cfg)
	if err != nil {
		log.Fatalf("Failed to decode env: %v", err)
	}
	log.Println("Env setup finished")
	return cfg
}
