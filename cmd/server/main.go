package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/olawolu/log-kini/internal/server"
)

func main() {
	loadEnvironment()
	port := os.Getenv("PORT")
	log.Printf("Starting server on port: %v", port)
	srv := server.NewHTTPServer(fmt.Sprintf(":%v", port))
	log.Fatal(srv.ListenAndServe())
}

func loadEnvironment()  {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}