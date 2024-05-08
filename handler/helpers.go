package handler

/*
helpers.go: part of the handler package. Contains helper functions to the handler functions.
*/

import (
	"assignment1/utils"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"slices"
	"strings"
	"sync"
)

/*
floatPrecision function that truncates floating point numbers to a fixed decimal
Found on: https://gosamples.dev/round-float/
*/
func floatPrecision(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

/*
AppendUniqueAuthors function that appends unique authors to the slice of authors
Takes a url string that it decodes, and for every book in that url, it appends
the author so long that it is unique
*/
func AppendUniqueAuthors(w http.ResponseWriter, u string, a *[]utils.Author, m *sync.Mutex, wg *sync.WaitGroup) {

	// Empty gutendex response instance
	response := utils.GutendexResponse{}
	// Call the GetAndDecode function to get the response from the URL, return if error
	if GetAndDecode(w, u, 503, &response) != nil {
		return
	}

	// We then append the unique authors, and lock the mutex around the uniqueness check
	for _, book := range response.Results {
		for _, author := range book.Authors {
			m.Lock()
			if !slices.Contains(*a, author) {
				*a = append(*a, author)
			}
			m.Unlock()
		}
	}
	// Go routine tells the waitgroup it is done
	wg.Done()
}

/*
FindCountryNamesAndCode function that finds the country names for the language from the url.
*/
func FindCountryNamesAndCode(w http.ResponseWriter, cc string) ([]string, []string) {

	// Return variables
	var countryNames []string
	var isoCodes []string

	// Utility variables
	var responses []utils.ReadershipResponse
	url := utils.LANGUAGE2COUNTRIES + cc + "/"

	if GetAndDecode(w, url, 503, &responses) != nil {
		return nil, nil
	}

	// We append the country names and iso codes to the string slice we return from the function
	for i := 0; i < len(responses); i++ {
		countryNames = append(countryNames, responses[i].Country)
		isoCodes = append(isoCodes, responses[i].IsoCode)
	}

	return countryNames, isoCodes
}

/*
FindNumOfReaders function that finds the number of potential readers given a country.
The functions uses the provided country to query the REST Countries API and find the population of the country.
If the country has more than one value it will sum the populations and returns this number.
*/
func FindNumOfReaders(w http.ResponseWriter, c string) int {

	// Make a slice of empty response instances
	var response []utils.PopulationResponse
	url := utils.COUNTRIES + "alpha/" + strings.ToLower(c) + "/"

	if GetAndDecode(w, url, 503, &response) != nil {
		return 0
	}

	// Variable to hold the sum of the populations
	numOfReaders := 0

	// Go over the slices decoded from the URL and sum the populations
	for _, name := range response {
		numOfReaders += name.Population
	}

	return numOfReaders
}

/*
GetAndDecode function that gets a URL and decodes the response to a destination struct
Portrays the appropriate error message depending on the status code given in the parameter code.
The parameter dest is the struct to which the response is decoded to.
The function returns the error if it is thrown.

The purpose of the function is to clean up the function body of the handler functions.
*/
func GetAndDecode(w http.ResponseWriter, url string, code int, dest interface{}) error {

	// GET request to the URL
	urlGet, err := http.Get(url)
	if err != nil {
		if code == 500 {
			http.Error(w, "Error: Internal service unavailable", http.StatusInternalServerError)
		} else if code == 503 {
			http.Error(w, "Error: External service unavailable", http.StatusServiceUnavailable)
		}
		log.Print("Error during geting of URL: " + url)
		return err
	}

	// Decode and store the response in the destination struct
	decoder := json.NewDecoder(urlGet.Body)
	err = decoder.Decode(dest)
	if err != nil {
		http.Error(w, "Error during decoding of response at "+url, http.StatusInternalServerError)
		log.Print("Error during decoding of response in GetAndDecode")
		return err
	}
	return nil
}
