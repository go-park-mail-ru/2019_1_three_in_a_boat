package http_utils

import "testing"

func TestNewReport(t *testing.T) {
	fields := []string{"foo", "bar", "spam", "eggs", "Pepe the Frog"}
	r := NewReport(fields...)

	if r.Ok {
		t.Fatal("report is ok by default (must not be)")
	}

	for _, field := range fields {
		if _, ok := r.Fields[field]; !ok {
			t.Errorf("missing %s in report", field)
		} else if r.Fields[field].Ok != false || r.Fields[field].Errors != nil {
			t.Errorf("field report for %s is incorrect", field)
		}
	}

}
