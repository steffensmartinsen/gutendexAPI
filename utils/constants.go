package utils

/*
constants.go: part of the utils package. Contains constants for APIs, URLs, and Version.
*/

// DEFAULT_PATH The root path.
const DEFAULT_PATH = "/"

// BOOKCOUNT_PATH The path to the bookcount endpoint.
const BOOKCOUNT_PATH = "/librarystats/" + VERSION + "/bookcount/"

// READERSHIP_PATH The path to the readership endpoint
const READERSHIP_PATH = "/librarystats/" + VERSION + "/readership/"

// STATUS_PATH The path to the status endpoint.
const STATUS_PATH = "/librarystats/" + VERSION + "/status/"

// LOCAL_ROOT The root of the localhost host and port 8080.
const LOCAL_ROOT = "http://localhost:8080"

// LOCAL_BOOKCOUNT The bookcount endpoint on the localhost host and port 8080.
const LOCAL_BOOKCOUNT = LOCAL_ROOT + BOOKCOUNT_PATH

// LOCAL_READERSHIP The readership endpoint on the localhost host and port 8080.
const LOCAL_READERSHIP = LOCAL_ROOT + READERSHIP_PATH

// LOCAL_STATUS The status endpoint on the localhost host and port 8080.
const LOCAL_STATUS = LOCAL_ROOT + STATUS_PATH

// GUTENDEX The url to the Gutendex API root.
const GUTENDEX = "http://129.241.150.113:8000/books/"

// LANGUAGE2COUNTRIES The url ot the language2countries API root.
const LANGUAGE2COUNTRIES = "http://129.241.150.113:3000/language2countries/"

// COUNTRIES The url to the countries API root.
const COUNTRIES = "http://129.241.150.113:8080/v3.1/"

// VERSION The current version of the service.
const VERSION = "v1"
