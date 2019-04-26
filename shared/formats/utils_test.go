package formats

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestMakeErrorResponse(t *testing.T) {
	r := http.Request{}

	bytes := MakeErrorResponse(&r, "foo")

	var data JSONResponseData
	err := json.Unmarshal(bytes, &data)

	if err != nil {
		t.Fatal("failed to unmarshal JSON")
	}

	if data.Message != "foo" || data.Status != StatusMap[false] || data.Data != nil {
		t.Error("JSON marshaled incorrectly")
	}
}

func TestMakeClientErrorResponse(t *testing.T) {
	r := http.Request{}

	bytes, err := MakeClientErrorResponse(&r, "bar", "foo")
	if err != nil {
		t.Fatal("failed to marshal JSON")
	}

	var data JSONResponseData
	err = json.Unmarshal(bytes, &data)

	if err != nil {
		t.Fatal("failed to unmarshal JSON")
	}

	if data.Message != "foo" || data.Status != StatusMap[false] || data.Data != "bar" {
		t.Error("JSON marshaled incorrectly")
	}
}

func TestMakeSuccessResponse(t *testing.T) {
	r := http.Request{}

	bytes, err := MakeSuccessResponse(&r, "bar")
	if err != nil {
		t.Fatal("failed to marshal JSON")
	}

	var data JSONResponseData
	err = json.Unmarshal(bytes, &data)

	if err != nil {
		t.Fatal("failed to unmarshal JSON")
	}

	if data.Message != "" || data.Status != StatusMap[true] || data.Data != "bar" {
		t.Error("JSON marshaled incorrectly")
	}
}
