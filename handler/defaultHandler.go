package handler

/*
deafultHandler.go: part of the handler package. Returns a landing page to the end user at the root level.
*/

import (
	"net/http"
)

/*
EmptyHandler as default handler function that handles requests that go to root
*/
func EmptyHandler(w http.ResponseWriter, r *http.Request) {

	// Ensure the client interprets header as HTML
	w.Header().Set("content-type", "text/html")

	// Serve the index.html file as landing page at root
	http.ServeFile(w, r, "html/index.html")

}
