package handler

/*
statusHandler.go: part of the handler package. Returns the current status code from the APIs that the
service relies on, and the amount of time in seconds the service has been up.
*/

import (
	"assignment1/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// StartTime: Global variable to store the start time of the service
var StartTime time.Time

/*
GetUptime function that is initiated in main when the service is started.
Created here to get access in statusHandler to return the information to the client.
*/
func GetUptime() {
	StartTime = time.Now()
}

/*
StatusHandler as a handler function that calls the correct function based on the http method
*/
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		statusGetRequest(w)
	} else {
		http.Error(w, "Unsupported request method '"+r.Method+"' . Only "+
			http.MethodGet+" is supported.", http.StatusNotImplemented)
		return
	}
}

/*
statusGetRequest that handles the GET request sent to the status endpoint
*/
func statusGetRequest(w http.ResponseWriter) {

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// We run GET on all the services we depend on
	gutendexApi, err := http.Get(utils.GUTENDEX)
	if err != nil {
		log.Fatal("Error during GET request to Gutendex API")
	}
	languageApi, err := http.Get(utils.LANGUAGE2COUNTRIES)
	if err != nil {
		log.Fatal("Error during GET request to Language2Countries API")
	}
	countriesApi, err := http.Get(utils.COUNTRIES + "all")
	if err != nil {
		log.Fatal("Error during GET request to Countries API")
	}

	// Create an empty output instance for encoding
	output := utils.Status{}

	// Store the status code of the APIs
	output.GutendexAPI = gutendexApi.StatusCode
	output.LanguageAPI = languageApi.StatusCode
	output.CountriesAPI = countriesApi.StatusCode
	output.Version = utils.VERSION

	// Convert the time to string and append a 's' to the end
	outputTime := floatPrecision(time.Since(StartTime).Seconds(), 1)
	output.Uptime = strconv.FormatFloat(outputTime, 'f', 1, 64) + "s"

	// Encode the status codes back to the client
	encoder := json.NewEncoder(w)
	err = encoder.Encode(output)
	if err != nil {
		http.Error(w, "Error during encoding of status messages", http.StatusInternalServerError)
		log.Print("Error during encoding of status messages")
		return
	}
}
