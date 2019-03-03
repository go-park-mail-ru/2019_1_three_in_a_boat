package db

import (
	"time"
)

var GenderMap = [...]string{ // genderMap lmao
	"Не указан",
	"Мужской",
	"Женский",
	"Другой",
}

const (
	tableName = "user"
)

type user struct {
	id         uint64
	HighScore  int
	Gender     int
	Username   string
	FirstName  string
	LastName   string
	Email      string
	SignupDate time.Time
}

func (u *user) GetPK() uint64 {
	return u.id
}

func (u *user) setPK(pk uint64)  {
	u.id = pk
}

func (u *user) GetTableName() string {
	return tableName
}

func (u *user) Save(database Database) (err error) {
	if u.id == 0 {
		err = database.Store(u)
	} else {
		err = database.Update(u)
	}
	return
}

func NewUser() *user {
	return &user{}
}
