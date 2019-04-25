package formats

import (
	"encoding/json"
	"net/http"
)

// Makes a JSONResponse-compatible error response based on the message.
// Since formats.Marshal only returns error on incompatible types and values,
// this function is guaranteed to never fail. Provided for convenience.
func MakeErrorResponse(r *http.Request, message string) []byte {
	response, _ := json.Marshal(JSONResponse{
		r,
		false,
		message,
		nil})
	return response
}

// Makes a JSONResponse-compatible error response based on the message.
// Since formats.Marshal only returns error on incompatible types and values,
// this function is guaranteed to never fail. Provided for convenience.
func MakeClientErrorResponse(r *http.Request, data interface{},
	message string) ([]byte, error) {
	return json.Marshal(JSONResponse{r, false, message, data})
}

// Makes a JSONResponse-compatible error response based on the message.
// Provided for convenience. Unlike MakerErrorResponse, might fail if data
// parameter can not be properly marshaled.
func MakeSuccessResponse(r *http.Request, data interface{}) ([]byte, error) {
	return json.Marshal(JSONResponse{r, true, "", data})
}
