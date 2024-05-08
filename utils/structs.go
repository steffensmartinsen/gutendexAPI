package utils

/*
structs.go: part of the utils package. Contains all structs used in the input and output from the APIs
and returned back to the user.
*/

// Author struct to hold the Author data
type Author struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
}

// Book struct to hold the data of the books
type Book struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Authors []Author `json:"authors,omitempty"`
}

// GutendexResponse to hold the desired data from the API
type GutendexResponse struct {
	Language string `json:"language"`
	Next     string `json:"next"`
	Count    int    `json:"count"`
	Results  []Book `json:"results"`
	Authors  int
}

// GutendexOutput struct to hold the desired data from the bookcount endpoint
type GutendexOutput struct {
	Language string  `json:"language"`
	Count    int     `json:"count"`
	Authors  int     `json:"authors"`
	Fraction float64 `json:"fraction"`
}

// ReadershipOutput struct to hold the desired data from the readership endpoint
type ReadershipOutput struct {
	Country   string `json:"country"`
	Isocode   string `json:"isocode"`
	Books     int    `json:"books"`
	Authors   int    `json:"authors"`
	Readerhip int    `json:"readership"`
}

// ReadershipResponse struct to decode the data from the countries API
type ReadershipResponse struct {
	Country string `json:"Official_Name"`
	IsoCode string `json:"ISO3166_1_Alpha_2"`
}

// PopulationResponse struct to decode the data from the REST countries API
type PopulationResponse struct {
	Population int `json:"population"`
}

// Status struct to hold the data of the status endpoint
type Status struct {
	GutendexAPI  int    `json:"gutendexapi"`
	LanguageAPI  int    `json:"languageapi"`
	CountriesAPI int    `json:"countriesapi"`
	Version      string `json:"version"`
	Uptime       string `json:"uptime"`
}
