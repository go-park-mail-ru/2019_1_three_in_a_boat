package chat_db

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"time"
)

type Message struct {
	Pk        int64        `json:"mid"`
	Uid       db.NullInt64 `json:"uid,omitempty"`
	Text      string       `json:"text"`
	Timestamp time.Time    `json:"timestamp"`
}

func (m *Message) Save(_db db.Queryable) error {
	if m.Pk == 0 {
		return _db.QueryRow(
			`INSERT INTO message (uid, message) VALUES ($1, $2) RETURNING ID`,
			m.Uid, m.Text).Scan(&m.Pk)
	} else {
		return errors.New("editing not implemented")
	}
}

func GetLastNMessages(_db db.Queryable, limit int, offset int) (*sql.Rows, error) {
	limitStr := db.MakeLimitString(limit)
	offsetStr := db.MakeOffsetString(offset)

	return _db.Query(
		`SELECT m."id", m."uid", m."message", m."created"
            FROM message m ORDER BY m."created" DESC ` + limitStr + offsetStr)
}

func GetNMessagesSince(_db db.Queryable, limit int, firstMsgId int) (*sql.Rows, error) {
	limitStr := db.MakeLimitString(limit)

	return _db.Query(
		`SELECT m."id", m."uid", m."message", m."created"
            FROM message m WHERE id < $1 ORDER BY m."id" DESC `+limitStr, firstMsgId)
}

func MessageFromRow(row db.Scanner) (*Message, error) {
	m := &Message{}
	err := row.Scan(&m.Pk, &m.Uid, &m.Text, &m.Timestamp)
	if err == nil {
		return m, nil
	} else {
		return nil, err
	}
}
