// Package defines formats returned by the server or used internally, and some
// utility functions to simplify the work with these formats.
package formats

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats/pb"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

// The struct outlines full JSON response structure. It should not be used
// directly, use JSONResponse instead.
type JSONResponseData struct {
	Status  string      `json:"status"`
	Version string      `json:"version"`
	Date    string      `json:"responseTime"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	User    *pb.Claims  `json:"user"`
}

// The struct outlines basic JSON response structure. Additional metadata fields
// will be added by the MarshalJSON method
type JSONResponse struct {
	Request *http.Request
	Status  bool
	Message string
	Data    interface{}
}

// Converts JSONResponse into JSONResponseData, adding the missing fields automatically
func (jr JSONResponse) MarshalJSON() ([]byte, error) {
	u, _ := AuthFromContext(jr.Request.Context()) // if !ok just leave nil

	return json.Marshal(&JSONResponseData{
		StatusMap[jr.Status],
		settings.Version,
		time.Now().Format(time.RFC3339),
		jr.Data,
		jr.Message,
		u,
	})
}

// different from the db.DateFormat! DB one is fore pretty printing, this one
// is for accepting values from the API.
const DateFormat = "2-1-2006"

// JSON values returned to the client, indicating whether the response was
// completed successfully. Is redundant, considering http status code, so
// provided just for convenience
var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}
