package formats

import "encoding/json"

// Makes a JSONResponse-compatible error response based on the message.
// Since formats.Marshal only returns error on incompatible types and values,
// this function is guaranteed to never fail. Provided for convenience.
func MakeErrorResponse(message string) []byte {
	response, _ := json.Marshal(JSONResponse{
		false,
		message,
		nil})
	return response
}

// Makes a JSONResponse-compatible error response based on the message.
// Provided for convenience. Unlike MakerErrorResponse, might fail if data
// parameter can not be properly marshaled.
func MakeSuccessResponse(data interface{}) ([]byte, error) {
	return json.Marshal(JSONResponse{true, "", data})
}
