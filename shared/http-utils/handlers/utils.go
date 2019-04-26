package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
)

// Handles a successful response: tries to marshal data, writes it to res.
// If an error occurs, calls handle5XXDepth with the diagnostic message
func handle2XXDepth(depth int, status int, w http.ResponseWriter,
	r *http.Request, data interface{}) {
	jsonResponse, err := formats.MakeSuccessResponse(r, data)

	if err != nil {
		handle5XXDepth(http.StatusInternalServerError, 1+depth, w, r,
			formats.ErrJSONMarshalFailure, err)
	} else {
		w.WriteHeader(status)
		writeSuccessResponse(1+depth, w, r, jsonResponse, strconv.Itoa(status)+" OK")
	}
}

// Handles a user error: tries to marshal data, writes it to res. If an error
// occurs, calls handle5XXDepth with the diagnostic message The only difference
// from the handle2XX depth is this function also logs the message and returns
// it to the client.
func handleInvalidDataDepth(depth int, status int, w http.ResponseWriter,
	r *http.Request, data interface{}, msg string) {
	jsonResponse, err := formats.MakeClientErrorResponse(r, data, msg)

	if err != nil {
		handle5XXDepth(http.StatusInternalServerError, 1+depth, w, r,
			formats.ErrJSONMarshalFailure, err)
	} else {
		w.WriteHeader(status)
		writeSuccessResponse(1+depth, w, r, jsonResponse, msg)
	}
}

// Handles client errors: logs msg and forwards it to writeErrorResponse
// Unlike HandleInvalidData, does not serialize an interface into JSON, just
// sends the error message with empty data
func handle4XXDepth(depth int, status int, w http.ResponseWriter,
	r *http.Request, clientMsg, msg string) {
	w.WriteHeader(status)
	writeErrorResponse(1+depth, w, r, clientMsg, msg)
}

// Handles server errors: formats msg, err and forwards it to the writeErrorResponse
func handle5XXDepth(depth int, status int, w http.ResponseWriter, r *http.Request,
	msg string, err error) {
	w.WriteHeader(status)
	writeErrorResponse(1+depth, w, r, msg,
		fmt.Sprintf("%s (err: %s)", msg, err))
}

// Used to write an error message. Logs and returns it to the client. If sending
// the error to the client fails, logs that too. Assumes headers have already
// been written.
func writeErrorResponse(depth int, w http.ResponseWriter, r *http.Request,
	clientMsg string, msg string) {
	var responseError error = nil
	_, responseError = w.Write(formats.MakeErrorResponse(r, clientMsg))

	if responseError != nil {
		LogError(1+depth, fmt.Sprintf("%s while processing %s: %s",
			formats.ErrResponseWriterFailure, msg, responseError), r)
	} else {
		LogError(1+depth, msg, r)
	}
}

// Used to write a success message. Logs and return response to the client
// if sending the response to the client fails, logs that too
func writeSuccessResponse(depth int, w http.ResponseWriter, r *http.Request,
	response []byte, msg string) {
	_, err := w.Write(response)
	if err != nil {
		LogError(1+depth, fmt.Sprintf("%s while processing %s: %s",
			formats.ErrResponseWriterFailure, msg, err), r)
	} else {
		LogInfo(1+depth, msg, r)
	}
}
