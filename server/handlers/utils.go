package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"net/http"
)

// Handles a successful response: tries to marshal data, writes it to res.
// If an error occurs, calls handle5XXDepth with the diagnostic message
func handle2XXDepth(depth int, status int, w http.ResponseWriter,
	r *http.Request, data interface{}, returnedBy string) {
	jsonResponse, err := formats.MakeSuccessResponse(data)

	if err != nil {
		handle5XXDepth(http.StatusInternalServerError, 1+depth, w, r,
			formats.ErrJSONMarshalFailure, returnedBy+"/MakeSuccessResponse", err)
	} else {
		w.WriteHeader(status)
		writeSuccessResponse(1+depth, w, r, jsonResponse, "200 OK")
	}
}

// Handles client errors: logs msg and forwards it to writeErrorResponse
func handle4XXDepth(depth int, status int, w http.ResponseWriter,
	r *http.Request, msg string) {
	w.WriteHeader(status)
	writeErrorResponse(1+depth, w, r, msg, msg)
}

// Handles server errors: formats msg, returnedBy, err and forwards it to the
// writeErrorResponse
func handle5XXDepth(depth int, status int, w http.ResponseWriter, r *http.Request,
	msg string, returnedBy string, err error) {
	w.WriteHeader(status)
	writeErrorResponse(1+depth, w, r, msg,
		fmt.Sprintf("%s (returned by %s as %s)", msg, returnedBy, err))
}

// used to write an error message to logs and return it to the client
// if sending the error to the client fails, logs that too
func writeErrorResponse(depth int, w http.ResponseWriter, r *http.Request,
	clientMsg string, msg string) {
	var responseError error = nil
	_, responseError = w.Write(formats.MakeErrorResponse(clientMsg))

	if responseError != nil {
		logError(1+depth, fmt.Sprintf("%s while processing %s: %s",
			formats.ErrResponseWriterFailure, msg, responseError), r)
	} else {
		logError(1+depth, msg, r)
	}
}

// used to write an success message to logs and return response to the client
// if sending the response to the client fails, logs that too
func writeSuccessResponse(depth int, w http.ResponseWriter, r *http.Request,
	response []byte, msg string) {
	_, err := w.Write(response)
	if err != nil {
		logError(1+depth, fmt.Sprintf("%s while processing %s: %s",
			formats.ErrResponseWriterFailure, msg, err), r)
	} else {
		logInfo(1, msg, r)
	}
}
