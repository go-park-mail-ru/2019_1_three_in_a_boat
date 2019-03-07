package formats

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"time"
)

// The struct outlines full JSON response structure. It should not be used directly,
// use JSONResponse instead.
type _JSONResponseData struct {
	Status  string      `formats:"status"`
	Version string      `formats:"version"`
	Date    string      `formats:"responseTime"`
	Data    interface{} `formats:"data"`
	Message string      `formats:"message,omitempty"`
}

// The struct outlines basic JSON response structure. Additional metadata fields
// will be added by the MarshalJSON method
type JSONResponse struct {
	Status  bool
	Message string
	Data    interface{}
}

// Converts JSONResponse into _JSONResponseData, adding the missing fields automatically
func (jr JSONResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&_JSONResponseData{
		settings.StatusMap[jr.Status],
		settings.Version,
		time.Now().Format(time.RFC3339),
		jr.Data,
		jr.Message,
	})
}

// Represents a single field in a validation result
type JSONFieldValidation struct {
	Errors    []string `formats:"errors"`
	FieldName string   `formats:"fieldName"`
}

// Represents a form in a validation result
type JSONFormValidation struct {
	Fields []JSONFieldValidation `formats:"fields"`
}
