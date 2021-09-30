package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic(fmt.Sprintf("Environment key %s not found, recheck your .env file.", key))
	}
	return value
}

func GetEnv() map[string]string {
	env := make(map[string]string)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env["CREDENTIALS"] = getEnv("CREDENTIALS")
	env["TOKEN"] = getEnv("TOKEN")
	env["SHEETID"] = getEnv("SHEETID")

	return env
}
