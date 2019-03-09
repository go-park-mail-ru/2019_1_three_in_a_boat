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

// Logs 405 and sends back an error message, if possible
func Handle405(w http.ResponseWriter, r *http.Request) {
	handle4XXDepth(1, http.StatusMethodNotAllowed, w, r,
		formats.Err405, formats.Err405)
}

// Logs 404 and sends back an error message, if possible
func Handle404(w http.ResponseWriter, r *http.Request) {
	handle4XXDepth(1, http.StatusNotFound, w, r,
		formats.Err404, formats.Err404)
}

// Logs 404 and sends back an error message, if possible
func Handle403(w http.ResponseWriter, r *http.Request) {
	handle4XXDepth(1, http.StatusForbidden, w, r,
		formats.Err403, formats.Err403)
}

// Handles request provided in an invalid format, e.g. invalid JSON. Should
// never happen if the app is working as intended and the user isn't trying to
// do something fishy
func Handle400(w http.ResponseWriter, r *http.Request, clientMsg, msg string) {
	handle4XXDepth(1, http.StatusBadRequest, w, r,
		clientMsg, msg)
}

// Handles user error - generally a form validation error
func HandleInvalidData(w http.ResponseWriter, r *http.Request,
	clientMsg interface{}, msg string, returnedBy string) {
	handleInvalidDataDepth(1, http.StatusBadRequest, w, r, clientMsg, msg, returnedBy)
}
