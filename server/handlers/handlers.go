// Package defines helper handlers for some common http statuses
package handlers

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"net/http"
)

// Processes response for sending: if it can't be marshaled, handles 500
// the same way Handle500 does. Otherwise writes the message and logs 200 OK.
func Handle200(w http.ResponseWriter, r *http.Request,
	response interface{}, returnedBy string) {
	handle2XXDepth(1, http.StatusOK, w, r, response, returnedBy)
}

// Same as Handle200, but 201.
func Handle201(w http.ResponseWriter, r *http.Request,
	response interface{}, returnedBy string) {
	handle2XXDepth(1, http.StatusCreated, w, r, response, returnedBy)
}

// Logs the error and sends back an error message, if possible
func Handle500(w http.ResponseWriter, r *http.Request,
	msg string, returnedBy string, err error) {
	handle5XXDepth(1, http.StatusInternalServerError, w, r, msg, returnedBy, err)
}

// Logs 404 and sends back an error message, if possible
func Handle405(w http.ResponseWriter, r *http.Request) {
	handle4XXDepth(1, http.StatusMethodNotAllowed, w, r, formats.Err405)
}

// Logs 404 and sends back an error message, if possible
func Handle404(w http.ResponseWriter, r *http.Request) {
	handle4XXDepth(1, http.StatusNotFound, w, r, formats.Err404)
}
