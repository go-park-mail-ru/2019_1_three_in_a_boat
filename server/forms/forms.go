package forms

const minPasswordStrnegth = 0 // 0..4
const dateFormat = "2-1-2006"
const emailExistsCheck = false

type FormReport struct {
	Ok     bool                   `json:"ok"`
	Fields map[string]FieldReport `json:"fields"`
}

func MakeReport() *FormReport {
	return &FormReport{false,
		map[string]FieldReport{
			"username": {},
			"password": {},
			"email":    {},
			"name":     {},
			"lastname": {},
			"date":     {},
		},
	}
}

type FieldReport struct {
	Ok     bool     `json:"ok"`
	Errors []string `json:"errors"`
}
