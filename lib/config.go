package lib

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	DB_HOST          string
	DB_NAME          string
	DB_USER          string
	DB_PASSWORD      string
	DB_PORT          string
	IMAGEKIT_API_KEY string
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	// Return the value of the variable
	env := os.Getenv(key)
	if env == "" {
		log.Fatalf("the %s key not found in env!", key)
	}
	return env
}

func init() {
	DB_HOST = Config("DB_HOST")
	DB_NAME = Config("DB_NAME")
	DB_USER = Config("DB_USER")
	DB_PASSWORD = Config("DB_PASSWORD")
	DB_PORT = Config("DB_PORT")
	IMAGEKIT_API_KEY = Config("IMAGEKIT_API_KEY")
}
