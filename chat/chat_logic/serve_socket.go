package chat_logic

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type ChatSocket struct {
	Conn *websocket.Conn
	WriteMu sync.Mutex
	Request *http.Request
	Id		  string
}

func NewChatSocket(conn *websocket.Conn, r *http.Request) *ChatSocket {

	return &ChatSocket{conn, sync.Mutex{}}
}

func (cs * ChatSocket) Read() bool {
	err :=
}

func (cs *ChatSocket) ReadLoop() {
	for cs.Read() {

	}
}

func (cs *ChatSocket) Run() {
	cs.WriteMu.Lock()
	defer cs.WriteMu.Unlock()
}

func (cs *ChatSocket) WriteJSON(v interface{}) {
	cs.WriteMu.Lock()
	defer cs.WriteMu.Unlock()
	err := cs.Conn.WriteJSON(v)
	if handlers.WSHandleErrForward(cs.Request, formats.ErrWebSocketFailure, cs.Id, err) != nil {
		_ = cs.Conn.Close()
	}
}


