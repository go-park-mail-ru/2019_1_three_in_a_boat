package forms

const minPasswordStrnegth = 0 // 0..4
const dateFormat = "2-1-2006"
const emailExistsCheck = false

type Report struct {
	Ok     bool                   `json:"ok"`
	Fields map[string]FieldReport `json:"fields"`
}

func NewReport(fields ...string) *Report {
	report := &Report{false,
		make(map[string]FieldReport, len(fields)),
	}

	for _, field := range fields {
		report.Fields[field] = FieldReport{false, nil}
	}
	return report
}

type FieldReport struct {
	Ok     bool     `json:"ok"`
	Errors []string `json:"errors"`
}
