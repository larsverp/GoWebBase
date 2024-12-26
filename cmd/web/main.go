package main

import (
	"example.com/go-web-base/cmd/web/handler"
	"example.com/go-web-base/internal/application"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	h := handler.BaseHandler{App: app}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		h.IndexPage(w, r)
	})

	fmt.Println("Starting server on :8082")
	err = http.ListenAndServe(":8082", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
