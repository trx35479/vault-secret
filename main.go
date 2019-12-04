package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/vault-secret/controller"
)

// Initialise Env - Temporary solution if run thru local environment
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// Main functions
func main() {
	logPath := "vault-secret.log"

	controller.LogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/secret", controller.Handler)

	err := http.ListenAndServe(":8080", controller.Logging(http.DefaultServeMux))

	if err != nil {
		log.Fatal(err)
	}
}
