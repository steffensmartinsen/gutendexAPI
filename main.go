package main

import (
	"assignment1/handler"
	"assignment1/utils"
	"log"
	"net/http"
	"os"
)

func main() {

	// Start the timer of the service
	handler.GetUptime()

	// Handle port assignment (locally overrides if it does not find a environment variable)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set, using default port 8080")
		port = "8080"
	}

	// Handler endpoints
	http.HandleFunc(utils.DEFAULT_PATH, handler.EmptyHandler)
	http.HandleFunc(utils.BOOKCOUNT_PATH, handler.BookcountHandler)
	http.HandleFunc(utils.STATUS_PATH, handler.StatusHandler)
	http.HandleFunc(utils.READERSHIP_PATH, handler.ReadershipHandler)

	// Starting server
	log.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
