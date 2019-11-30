package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "SRS server: ", log.LstdFlags)
}

func main() {
	/**
	 * Load configuration variables from ".env" file in directory
	 * or from Bash Environment.
	**/
	logger.Println("Initializing Environment.")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error accessing .env file.")
	}

	/**
	 * Connect to databases with loaded config.
	**/
	ConnectDB()
	ConnectRedis()

	r := createRouter()
	http.Handle("/", r)
	logger.Println("Server starting on 0.0.0.0:8080...")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
