// Package defines helper handlers for some common http usages
package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
)

// Processes response for sending: if it can't be marshaled, handles 500
// the same way Handle500 does. Otherwise writes the message and logs 200 OK.
func Handle200(w http.ResponseWriter, r *http.Request,
	response interface{}) {
	handle2XXDepth(1, http.StatusOK, w, r, response)
}

// Same as Handle200, but 201.
func Handle201(w http.ResponseWriter, r *http.Request,
	response interface{}) {
	handle2XXDepth(1, http.StatusCreated, w, r, response)
}

// Logs the error and sends back an error message, if possible
func Handle500(w http.ResponseWriter, r *http.Request, msg string, err error) {
	handle5XXDepth(1, http.StatusInternalServerError, w, r, msg, err)
}

// If err != nil, calls handle 500. Returns the original err eitherway, so it
// can be used in the context if HandleErrForward(...) != nil { return }
func HandleErrForward(w http.ResponseWriter, r *http.Request,
	msg string, err error) error {
	if err == nil {
		return nil
	} else {
		handle5XXDepth(1, http.StatusInternalServerError, w, r, msg, err)
		return err
	}
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

// Logs 403 and sends back an error message, if possible
func Handle403(w http.ResponseWriter, r *http.Request) {
	handle4XXDepth(1, http.StatusForbidden, w, r,
		formats.Err403, formats.Err403)
}

func Handle403Msg(w http.ResponseWriter, r *http.Request, msg string) {
	handle4XXDepth(1, http.StatusForbidden, w, r,
		msg, msg)
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
	clientMsg interface{}, msg string) {
	handleInvalidDataDepth(1, http.StatusBadRequest, w, r, clientMsg, msg)
}

// If !report.Ok, callse HandleInvalidData. Returns report eitherway, so it can
// be used in the context if !HandleReportForward(...).Ok { return }
func HandleReportForward(w http.ResponseWriter, r *http.Request,
	report *forms.Report) *forms.Report {
	if !report.Ok {
		handleInvalidDataDepth(1, http.StatusBadRequest,
			w, r, report, formats.ErrValidation)
	}
	return report
}

func WSLogInfo(r *http.Request, msg, connId string) {
	LogInfo(1, fmt.Sprintf("Connection %s: %s", connId, msg), r)
}

func WSLogError(r *http.Request, msg, connId string, err error) {
	LogError(1, fmt.Sprintf(
		"Connection %s: %s (%s)", connId, msg, err.Error()), r)
}

func WSHandleErrForward(r *http.Request, msg, connId string, err error) error {
	if err == nil {
		return nil
	} else {
		LogError(1, fmt.Sprintf(
			"Connection %s: %s (%s)", connId, msg, err.Error()), r)
		return err
	}
}
