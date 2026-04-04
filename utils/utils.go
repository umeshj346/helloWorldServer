package utils

import (
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load("D:\\Golang\\Backend\\http-server\\hello-world\\.env")
	if err != nil {
		return err
	}
	return nil
}