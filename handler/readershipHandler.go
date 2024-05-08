package handler

/*
readershipHandler.go: part of the handler package. Handles incoming URLs to the /readership endpoint.
The handler only deals with GET requests, and decodes the incoming URL to retrieve the number of books,
the number of unique authors for a given language, and retrieves the amount of countries that speak
the given language, and the amount of potential readers in those countries.
The handler can deal with an optional parameter which limits the amount of countries returned for the language.
*/

import (
	"assignment1/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/*
ReadershipHandler as a handler function that calls the correct function based on the http method
*/
func ReadershipHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		readershipGetRequest(w, r)
	} else {
		http.Error(w, "Unsupported request method '"+r.Method+"' . Only "+
			http.MethodGet+" is supported.", http.StatusNotImplemented)
		return
	}
}

/*
readershipGetRequest that handles the GET request sent to the readership endpoint
*/
func readershipGetRequest(w http.ResponseWriter, r *http.Request) {

	// Parse the URL
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Error during reading of URL", http.StatusBadRequest)
		log.Print("Error during parsing of URL")
		return
	}

	// Get the query parameters
	queryParams := u.Query()
	limitStr := queryParams.Get("limit")

	// Set limit to 0 if the parameter is empty, else convert the value to an int and store it in limit
	limit := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Error during generation of response", http.StatusInternalServerError)
			log.Print("Error during parsing of limit")
			return
		}
	}

	// Get the path
	path := u.Path

	// Append a / if the URL is lacking the final one
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// Split the path on the /
	pathSplit := strings.Split(path, "/")

	// Throw an error if the path beyond root is empty
	if len(pathSplit) == 5 {
		http.Error(w, "Error: Please provide a country code", http.StatusBadRequest)
		return
	}

	// Throw an error if path is incomplete
	if len(pathSplit) < 5 {
		http.Error(w, "Error: Incomplete URL", http.StatusBadRequest)
		return
	}

	// Get the country code which will be on the 5th index if we are not at root
	countryCode := pathSplit[4]

	// String slices to store the country names and ISO codes, get the values from the function
	var countries []string
	var isoCodes []string
	countries, isoCodes = FindCountryNamesAndCode(w, countryCode)

	// Query our own bookcount endpoint to get the relevant values
	bookcountUrl := utils.LOCAL_BOOKCOUNT + "?language=" + countryCode

	// Temporary variable used for decoding
	var responses []utils.GutendexOutput

	if GetAndDecode(w, bookcountUrl, 500, &responses) != nil {
		// Return if error
		return
	}

	// If limit is more than amount of countries, we reduce limit to the amount of countries
	if limit > len(countries) || limit == 0 {
		limit = len(countries)
	}

	// Slice to hold the responses for encoding, and one single instance
	var outputs []utils.ReadershipOutput
	output := utils.ReadershipOutput{}

	// Iterate over the countries limit number of times
	for i := 0; i < limit; i++ {

		// Popularize the data members of the output struct
		output.Country = countries[i]
		output.Authors = responses[0].Authors
		output.Books = responses[0].Count
		output.Isocode = isoCodes[i]
		output.Readerhip = FindNumOfReaders(w, isoCodes[i])

		// Append the output struct to the outputs slice
		outputs = append(outputs, output)
	}

	// Sett the content-type to JSON
	w.Header().Add("Content-Type", "application/json")

	// Show output for every country that speaks the language (within the limit)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(outputs)
	if err != nil {
		http.Error(w, "Error during generation of response", http.StatusInternalServerError)
		log.Print("Error during encoding of structs.")
		return
	}

}
