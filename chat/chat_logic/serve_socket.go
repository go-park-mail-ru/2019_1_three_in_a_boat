package chat_logic

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/chat/chat_db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/chat"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type ChatSocket struct {
	Conn *websocket.Conn
	// only locks and unlocks in WriteJson, to avoid
	// deadlock/undefined behavior/panic/whatever, it must not be locked anywhere else
	WriteMu sync.Mutex
	Request *http.Request
	Id      string
}

type ServiceMessage struct {
	Code string `json:"message"`
	Type string `json:"type"`
}

func NewChatSocket(conn *websocket.Conn, r *http.Request) *ChatSocket {
	id := uuid.New().String()
	cs := &ChatSocket{conn, sync.Mutex{}, r, id}
	MainChat.Sockets.Store(id, cs)
	return cs
}

func (cs *ChatSocket) Close() {
	_ = cs.Conn.Close()
	MainChat.Sockets.Delete(cs.Id)
}

func (cs *ChatSocket) ReadJSON(v interface{}) bool {
	err := cs.Conn.ReadJSON(v)
	if handlers.WSHandleErrForward(cs.Request, formats.ErrWebSocketFailure, cs.Id, err) != nil {
		cs.Close()
		return false
	} else {
		return true
	}
}

func (cs *ChatSocket) WriteJSON(v interface{}) bool {
	cs.WriteMu.Lock()
	defer cs.WriteMu.Unlock()
	err := cs.Conn.WriteJSON(v)
	if handlers.WSHandleErrForward(cs.Request, formats.ErrWebSocketFailure, cs.Id, err) != nil {
		cs.Close()
		return false
	}
	return true
}

func (cs *ChatSocket) Read() bool {
	m := chat_db.Message{}
	if cs.ReadJSON(&m) {
		u, ok := formats.AuthFromContext(cs.Request.Context())
		m.Uid = 0
		if u != nil && ok {
			m.Uid = u.Uid
		}
		m.Timestamp = time.Now()

		err := m.Save(chat_settings.DB())
		if err != nil {
			handlers.WSLogError(cs.Request, formats.ErrSqlFailure, cs.Id, err)
			return cs.WriteJSON(ServiceMessage{formats.ErrSqlFailure, "error"})
		}
		WriteToAll(m)
		return true
	}

	return false
}

func (cs *ChatSocket) ReadLoop() {
	for cs.Read() {
	}
}

func (cs *ChatSocket) Run() {
	cs.WriteJSON(ServiceMessage{cs.Id, "id"})
	cs.WriteJSON(GetLastMessages())
	cs.ReadLoop()
}
