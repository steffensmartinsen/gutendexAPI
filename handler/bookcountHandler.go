package handler

/*
bookcountHandler.go: part of the handler package. Handles incoming URLs to the /bookcount endpoint.
The handler only deals with GET requests, and decodes the URL parameters to retrieve the amount of books,
the amount of unique authors, and the fraction of the total amount of books from the Gutenberg Library API,
in addition to providing the language from the URL.
*/

import (
	"assignment1/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

/*
BookcountHandler as a handler function that calls the correct function based on the
HTTPS method.
*/
func BookcountHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		bookCountGetRequest(w, r)
	} else {
		http.Error(w, "Unsupported request method '"+r.Method+"' . Only "+
			http.MethodGet+" is supported.", http.StatusNotImplemented)
		return
	}
}

/*
bookCountGetRequest function that handles GET request sent to the bookcount endpoint
*/
func bookCountGetRequest(w http.ResponseWriter, r *http.Request) {

	// Retrieve the count of ALL the books in the Gutenberg API
	// Can not make it a const in case it changes
	allBooks := utils.GutendexResponse{}

	// Call Get and Decode on the main page of the Gutendex API
	if GetAndDecode(w, utils.GUTENDEX, 503, &allBooks) != nil {
		// Return if error
		return
	}

	// Parse the URL
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Error during reading of URL", http.StatusBadRequest)
		log.Print("Error during parsing of URL")
		return
	}

	// Get the query parameters and language value from the parameter(s)
	queryParams := u.Query()
	language := queryParams.Get("language")

	// If queryParams == 0 we throw an error
	if len(queryParams) == 0 || language == "" {
		http.Error(w, "Error: Please provide a language", http.StatusBadRequest)
		return
	}

	// Split the country codes on the comma
	countryCodes := strings.Split(language, ",")

	// Slice of ouput structs to store the output for each countryCode
	var outputs []utils.GutendexOutput

	// For every country code we perform the get to the API to get the desired output
	for _, countryCode := range countryCodes {

		// Add the parameter to the url - slice the countryCode to make nor, norge etc work
		gutenUrl := utils.GUTENDEX + "?languages=" + countryCode[:2]

		// Variable for decoding and calling Get and Decode on the gutendex API
		gutendexResponse := utils.GutendexResponse{}
		if GetAndDecode(w, gutenUrl, 503, &gutendexResponse) != nil {
			return
		}

		// Slice that contains the unique authors for the current language
		var authors []utils.Author

		// Create empty instance of the output struct
		output := utils.GutendexOutput{}

		// Fill the output struct with the desired information
		output.Language = countryCode[:2] // To secure that only the two-letter ISO code is used
		output.Count = gutendexResponse.Count
		fraction := float64(gutendexResponse.Count) / float64(allBooks.Count)
		output.Fraction = floatPrecision(fraction, 4)

		// Store the amount of pages to be queried (there are a maximum of 32 books per page)
		pagesAmount := gutendexResponse.Count / 32

		// Make a waitgroup and a mutex lock for the go routines
		var wg sync.WaitGroup
		m := sync.Mutex{}

		// Add amount of pages + the default page to the waitgroup
		wg.Add(pagesAmount + 1)

		/*
			Make a go routine for the first page, and then for the rest of the pages
			Trying to make large queries faster, this reduced the runtime for finland from
			~46s to ~30s on my local network.
		*/
		go AppendUniqueAuthors(w, gutenUrl, &authors, &m, &wg)
		for i := 0; i < pagesAmount; i++ {
			// The first URL with a page parameter starts at page=2
			gutenUrl = utils.GUTENDEX + "?languages=" + countryCode[:2] + "&page=" + strconv.Itoa(i+2)
			go AppendUniqueAuthors(w, gutenUrl, &authors, &m, &wg)
		}

		// Wait for all the go routines to finish
		wg.Wait()
		output.Authors = len(authors)

		// Append the output struct to the outputs slice
		outputs = append(outputs, output)

	}

	// Set the header to JSON and Encode the gutendexResponses back to the client
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(outputs)
	if err != nil {
		http.Error(w, "Error during generation of response", http.StatusInternalServerError)
		log.Print("DEBUGGING: Error during encoding of structs.")
		return
	}

}
