package forms

import (
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
)

func TestCheckLength(t *testing.T) {
	cases := []struct {
		str string
		min int
		max int
		ok  bool
	}{
		{
			str: "foobar",
			min: 6,
			max: 10,
			ok:  true,
		},
		{
			str: "foobar",
			min: 8,
			max: 10,
			ok:  false,
		},
		{
			str: "foobarfoobar",
			min: 6,
			max: 10,
			ok:  false,
		},
		{
			str: "barbar",
			min: 0,
			max: 0,
			ok:  false,
		},
	}

	for _, c := range cases {
		if ok := CheckLength(c.str, c.min, c.max).Ok; ok != c.ok {
			t.Errorf(
				"%d <= %s <= %d: expceted %t, got %t", c.min, c.str, c.max, c.ok, ok)
		}
	}
}

func TestCheckLengthOptional(t *testing.T) {
	cases := []struct {
		str db.NullString
		min int
		max int
		ok  bool
	}{
		{
			str: db.NullString{"foobar", true},
			min: 6,
			max: 10,
			ok:  true,
		},
		{
			str: db.NullString{"foobar", true},
			min: 8,
			max: 10,
			ok:  false,
		},
		{
			str: db.NullString{"foobarfoobar", true},
			min: 6,
			max: 10,
			ok:  false,
		},
		{
			str: db.NullString{"", true},
			min: 5,
			max: 12,
			ok:  true,
		},
		{
			str: db.NullString{"", false},
			min: 5,
			max: 12,
			ok:  true,
		},
	}

	for _, c := range cases {
		if ok := CheckLengthOptional(&c.str, c.min, c.max).Ok; ok != c.ok {
			t.Errorf(
				"%d <= %s <= %d: expceted %t, got %t", c.min, c.str.String, c.max, c.ok, ok)
		}
	}
}
