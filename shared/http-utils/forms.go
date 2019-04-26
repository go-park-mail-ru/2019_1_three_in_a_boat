// The package provides tools for validating forms: serializing them, checking
// validation and creating DB objects based on the received data. Forms are
// designed in a way that does not make any changes - only validate the data
// and suggest changes. E.g., saving an image to the filesystem, saving an
// entity in the database has to be done outside of the forms.
package http_utils

import "errors"

// The error returned when a form is required to make an action that requires
// validation, but has not yet been validated. Note that the call to Validate()
// must be explicit since it returns a valuable and potentially time-consuming
// report
var ErrFormInvalid = errors.New("can not create/edit DB objects from an invalid form")

// The type returned by *Form.Validate. Ok is a global AND on all the Field.Ok
// of the report.
type Report struct {
	Ok     bool                   `json:"ok"`
	Fields map[string]FieldReport `json:"fields"`
}

// The type included in Report: note that errors might be nil, so for checking
// if field validation succeeded, check Ok first.
type FieldReport struct {
	Ok     bool     `json:"ok"`
	Errors []string `json:"errors"`
}

// Generates a report with given fields, provided for convenience. Ok for the
// report and all fields is set to false, error slices are set to nil.
func NewReport(fields ...string) *Report {
	report := &Report{false,
		make(map[string]FieldReport, len(fields)),
	}

	for _, field := range fields {
		report.Fields[field] = FieldReport{false, nil}
	}
	return report
}
