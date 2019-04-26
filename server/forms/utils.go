package forms

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	http_utils "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
)

// returns a fieldReport indicating min < len(str) < max. Length counted in
// runes, not in bytes
func CheckLength(str string, min, max int) http_utils.FieldReport {
	length := len([]rune(str))
	if length < min {
		return http_utils.FieldReport{false, []string{formats.ErrFieldTooShort}}
	}
	if length > max {
		return http_utils.FieldReport{false, []string{formats.ErrFieldTooLong}}
	}

	return http_utils.FieldReport{true, nil}
}

// returns a fieldReport indicating min < len(str) < max. Length counted in
// runes, not in bytes. Also allows empty strings, even if min is > 0.
func CheckLengthOptional(str *db.NullString, min, max int) http_utils.FieldReport {
	length := len([]rune(str.String))
	if length < min {
		if length != 0 {
			return http_utils.FieldReport{false, []string{formats.ErrFieldTooShort}}
		} else {
			return http_utils.FieldReport{true, []string{}}
		}
	}
	if length > max {
		return http_utils.FieldReport{false, []string{formats.ErrFieldTooLong}}
	}

	str.Valid = true
	return http_utils.FieldReport{true, []string{}}
}
