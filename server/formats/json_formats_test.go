package formats

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestJSONResponse_MarshalJSON(t *testing.T) {
	res := JSONResponse{
		Request: &http.Request{},
		Status:  true,
		Message: "foo",
		Data:    "bar",
	}
	bytes, err := res.MarshalJSON()
	if err != nil {
		t.Fatal("failed to marshal JSON")
	}

	var data JSONResponseData
	err = json.Unmarshal(bytes, &data)

	if err != nil {
		t.Fatal("failed to unmarshal JSON")
	}

	if data.Message != "foo" || data.Status != StatusMap[true] || data.Data != "bar" {
		t.Error("JSON marshaled incorrectly")
	}

}
