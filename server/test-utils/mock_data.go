package test_utils

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"time"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type Mock struct {
	*User
	*Author
}

var mockData []Mock

func GetMockData() []Mock {
	return mockData
}

// creates a User + Author, panics on any error, should only be used for mocking data
func makeMockUser(_db *sql.DB, username, email, password, firstName,
	lastName string, highScore int64, gender, img string, birthDate time.Time,
	devInfo string, description string) *User {
	a, err := NewAccount(username, email, password)
	panicOnError(err)
	p, err := NewProfile(0, NullString{firstName, true},
		NullString{lastName, true}, NullInt64{highScore, true},
		NullString{gender, true}, NullString{img, true},
		NullTime{birthDate, true},
	)
	panicOnError(err)
	u, err := NewUser(a, p)
	panicOnError(err)
	err, txError := u.Save(_db)
	panicOnError(txError)
	panicOnError(err)
	au := NewAuthor(u.Account.Pk, NullString{devInfo, true},
		NullString{description, true})
	err = au.Save(_db)
	panicOnError(err)
	mockData = append(mockData, Mock{u, au})
	return u
}

// used to create garbage data for debugging purposes - useless otherwise
func StoreMockData(_db *sql.DB) []*User {
	users := make([]*User, 3)

	users[0] = makeMockUser(_db, "alfaix", "example@exam.ple", "12345",
		"Арсен", "Китов", 10000, "male", "pepe.jpg",
		time.Date(1998, 16, 12, 0, 0, 0, 0, time.UTC),
		"фронт/бэк/фуллкек", "На самом деле он ничего не делал",
	)

	users[1] = makeMockUser(_db, "Kotyarich", "example2@exam.ple", "12345",
		"Никита", "Котов", 10000, "male", "nikita.jpg",
		time.Date(1973, 5, 2, 0, 0, 0, 0, time.UTC),
		"фронт", "Он тоже ничего не делал",
	)

	users[2] = makeMockUser(_db, "Rowbotman", "example3@exam.ple", "12345",
		"Андрей", "Прокопенко", 10000, "male", "pepe.jpg",
		time.Date(1989, 1, 11, 0, 0, 0, 0, time.UTC),
		"фронт", "И он тоже. Как это вообще работает?",
	)

	return users
}
