package routes

import (
	"encoding/json"
	"time"
)

// The struct outlines full JSON response structure. It should not be used directly,
// use JSONResponse instead.
type _JSONResponseData struct {
	Status  string      `json:"status"`
	Version string      `json:"version"`
	Date    string      `json:"responseTime"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// The struct outlines basic JSON response structure. Additional metadata fields
// will be added by the MarshalJSON method
type JSONResponse struct {
	Status  bool
	Message string
	Data    interface{}
}

// Makes a JSONResponse-compatible error response based on the message.
// Provided for convenience.
func MakeErrorResponse(message string) []byte {
	response, _ := json.Marshal(JSONResponse{false, message, nil})
	return response
}

// Makes a JSONResponse-compatible error response based on the message.
// Provided for convenience.
func MakeSuccessResponse(data interface{}) ([]byte, error) {
	return json.Marshal(JSONResponse{true, "", data})
}

// Converts JSONResponse into _JSONResponseData, adding the missing fields automatically
func (jr JSONResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&_JSONResponseData{
		StatusMap[jr.Status],
		Version,
		time.Now().Format(time.RFC3339),
		jr.Data,
		jr.Message,
	})
}

// Represents a single field in a validation result
type JSONFieldValidation struct {
	Errors    []string `json:"errors"`
	FieldName string   `json:"fieldName"`
}

// Represents a form in a validation result
type JSONFormValidation struct {
	Fields []JSONFieldValidation `json:"fields"`
}
