package main

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/chat/chat_logic"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

// this is actually a CHAD handler and he is very alpha
type ChatHandler struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // PepegaChamp
		origin := r.Header.Get("Origin")
		_, allowed := settings.GetAllowedOrigins()[origin]
		return allowed
	},
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if HandleErrForward(w, r, formats.ErrWebSocketFailure, err) != nil {
		return
	}

	chat := chat_logic.NewChatSocket(conn, r)
	chat.Run()
}

func (h *ChatHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
