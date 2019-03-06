package routes

import (
	"fmt"
	"github.com/google/logger"
	"net/http"
)

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
}

func setDefaultHeaders(res http.ResponseWriter) {
	for k, v := range defaultHeaders {
		res.Header().Set(k, v)
	}
}

// Calls handle500Depth with depth = 0. Provided for convenience
func handle500(h Handler, msg string, returnedBy string,
	err error, res http.ResponseWriter) {
	handle500Depth(1, h, msg, returnedBy, err, res)
}

// Handles an internal server error: tries to write 500 and msg to the client,
// logs failure (and failure to write the error message, if it happens)
// forwards err back to the caller for convenience. Depth is incremented and
// forwarded to the logger
func handle500Depth(depth int, h Handler, msg string, returnedBy string,
	err error, res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)

	var responseError error = nil
	if msg != "E_RESPONSE_WRITER_FAILURE" {
		// no point in trying to write again and pollute the log: it's likely a
		// premature disconnect
		_, responseError = res.Write(MakeErrorResponse(msg))
	}

	if responseError != nil {
		logger.ErrorDepth(1+depth, fmt.Sprintf(
			"%s: failed to write an error response %s (returned by %s as %s): %s\n",
			h.GetRoute(), msg, returnedBy, err, responseError))
	} else {
		logger.ErrorDepth(1+depth, fmt.Sprintf(
			h.GetRoute(), "%s: %s (returned by %s as %s)\n", msg, returnedBy, err))
	}

}

// Handles a successful response: tries to marshal data, writes it to res.
// Logs and returns an error if one occurs during marshaling or writing response
// logs Info if no error happened
func handle200(h Handler, data interface{}, dataSrc string,
	res http.ResponseWriter) {
	jsonResponse, err := MakeSuccessResponse(data)
	if err != nil {
		handle500Depth(1, h, "E_JSON_MARSHAL_FAILURE",
			"MakeSuccessResponse/"+dataSrc, err, res)
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write(jsonResponse)
	if err != nil {
		handle500Depth(1, h, "E_RESPONSE_WRITER_FAILURE",
			"res.Write", err, res)
	} else {
		logger.Infof("%s: 200 OK\n", h.GetRoute())
	}
}
