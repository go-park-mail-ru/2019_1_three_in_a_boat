package chat_logic

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/chat/chat_db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
	"github.com/google/logger"
	"sync"
)

type Chat struct {
	Sockets sync.Map
}

var MainChat = Chat{}

func WriteToAll(message chat_db.Message) {
	MainChat.Sockets.Range(func(_sockId, _sock interface{}) bool {
		sock := _sock.(*ChatSocket)
		sock.WriteJSON(message)
		return true
	})
}

func GetLastMessages() []*chat_db.Message {
	rows, err := chat_db.GetLastNMessages(settings.DB(), 20, 0)
	if err != nil {
		logger.Error(formats.ErrSqlFailure, ": ", err)
		return nil
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	var ret []*chat_db.Message
	for rows.Next() {
		m, err := chat_db.MessageFromRow(rows)
		if err != nil {
			logger.Error(formats.ErrSqlFailure, err)
			// still finish the loop and return what we can
		} else {
			ret = append(ret, m)
		}
	}

	return ret
}
