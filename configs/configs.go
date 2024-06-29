package configs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	err := LoadEnv()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
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

func LoadEnv() error {
	// open .env file
	file, err := os.Open(".env")
	if err != nil {
		return err
	}

	defer file.Close()

	// read the file into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return err
	}

	// unmarshal the buffer into a map
	replacedBuf := bytes.Replace(buf.Bytes(), []byte("\r\n"), []byte("\n"), -1)
	lines := strings.Split(string(replacedBuf), "\n")

	envMap := make(map[string]string)

	// read the key value pairs
	for _, line := range lines {
		values := strings.Split(line, "=")
		envMap[values[0]] = values[1]
	}

	// add the .env vars into the os env
	for key, value := range envMap {
		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
