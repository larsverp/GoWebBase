package main

import (
	"example.com/go-web-base/internal/application"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if os.Getenv("ENV") == "" {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db, err := application.NewDatabase()
	if err != nil {
		panic("cannot get database connection: " + err.Error())
	}

	app := application.Application{
		DB:  db,
		Log: application.PrintLnLogger{},
	}
}
