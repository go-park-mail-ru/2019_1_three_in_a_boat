package forms

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
)

// returns a fieldReport indicating min < len(str) < max. Length counted in
// runes, not in bytes
func checkLength(str string, min, max int) FieldReport {
	length := len([]rune(str))
	if length < min {
		return FieldReport{false, []string{formats.ErrFieldTooShort}}
	}
	if length > max {
		return FieldReport{false, []string{formats.ErrFieldTooLong}}
	}

	return FieldReport{true, nil}
}

// returns a fieldReport indicating min < len(str) < max. Length counted in
// runes, not in bytes. Also allows empty strings, even if min is > 0.
func checkLengthOptional(str *db.NullString, min, max int) FieldReport {
	length := len([]rune(str.String))
	if length < min {
		if length != 0 {
			return FieldReport{false, []string{formats.ErrFieldTooShort}}
		} else {
			return FieldReport{true, []string{}}
		}
	}
	if length > max {
		return FieldReport{false, []string{formats.ErrFieldTooLong}}
	}

	str.Valid = true
	return FieldReport{true, []string{}}
}
