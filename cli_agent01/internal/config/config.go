package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Model    string
	ApiUrl   string
	ApiToken string
	Timeout  int
	Retries  int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	retries, err := strconv.Atoi(os.Getenv("RETRIES"))
	if err != nil {
		return nil, err
	}

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Model:    os.Getenv("MODEL"),
		ApiUrl:   os.Getenv("API_URL"),
		ApiToken: os.Getenv("API_TOKEN"),
		Timeout:  timeout,
		Retries:  retries,
	}, nil
}
