package game

import (
	"encoding/json"
	"testing"
)

func TestInput_UnmarshalJSON(t *testing.T) {
	i := NewInput(0)
	err := json.Unmarshal([]byte(`{"angle":0.3}`), &i)
	if err != nil {
		t.Fatal(err.Error())
	}
}
