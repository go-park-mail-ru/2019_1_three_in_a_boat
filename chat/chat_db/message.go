package messages

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"time"
)

type Message struct {
	Pk        int64     `json:"mid,omitempty"`
	Uid       int64     `json:"uid,omitempty"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *Message) Save(_db db.Queryable) error {
	if m.Pk == 0 {
		return _db.QueryRow(
			`INSERT INTO message (uid, text) VALUES ($1, $2)`,
			m.Uid, m.Text).Scan(&m.Pk)
	} else {
		return errors.New("editing not implemented")
	}
}

func GetLastNMessages(_db db.Queryable, limit int, offset int) (*sql.Rows, error) {
	limitStr := db.MakeLimitString(limit)
	offsetStr := db.MakeOffsetString(offset)

	return _db.Query(
		`SELECT m."id", m."uid", m."text", m."timestamp",
            FROM message m ORDER BY m."timestamp" DESC ` + limitStr + offsetStr)
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
