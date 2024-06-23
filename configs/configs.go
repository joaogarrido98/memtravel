package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	JWTSecret  string
	JWTIssuer  string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return Config{
		Port:       os.Getenv("PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBAddress:  fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		JWTIssuer:  os.Getenv("JWT_ISSUER"),
	}
}
