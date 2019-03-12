package test_utils

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"github.com/google/logger"
	"io/ioutil"
)

func SetUp() {
		logger.Init("", false, false, ioutil.Discard)
		settings.SetDbParams("", "", "", "hexagon_test")
		// teardown is needed if the test stops halfway through: e.g., if stopped
		// by a debugger. Otherwise, main calls TearDown and there's no need to call
		// it here. Anyway, it's idempotent so it doesn't hurt
		TearDown()
		StoreMockData(settings.DB())
}

func TearDown() {
	// foreign keys first, or the constraint will complain
	_, err := settings.DB().Exec("DELETE FROM author;")
	if err != nil {
		panic(err)
	}
	_, err = settings.DB().Exec("DELETE FROM profile;")
	if err != nil {
		panic(err)
	}
	_, err = settings.DB().Exec("DELETE FROM account;")
	if err != nil {
		panic(err)
	}
}

func UserDataEqual(u1 db.UserData, u2 db.UserData) bool {
	if u1.SignupDate.Sub(u2.SignupDate).Seconds() < 1 {
		// json omits sub-second times, so this is a workaround
		// BirthData does not store time so no need to do it there
		u1.SignupDate = u2.SignupDate
	}
	return u1 == u2
}

func MockToUserData(m *Mock) db.UserData {
	udBytes, _ := json.Marshal(m.User)
	var userData db.UserData
	_ = json.Unmarshal(udBytes, &userData)

	return userData
}

func UserToUserData(u *db.User) db.UserData {
	var udBytes []byte
	if u != nil && u.Profile != nil && u.Account != nil {
		udBytes, _ = json.Marshal(u)
	}
	var userData db.UserData
	_ = json.Unmarshal(udBytes, &userData)

	return userData
}
