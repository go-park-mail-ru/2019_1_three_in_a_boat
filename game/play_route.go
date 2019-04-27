package main

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/game/game_logic"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

type SinglePlayHandler struct{}

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

func (h *SinglePlayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if HandleErrForward(w, r, formats.ErrWebSocketFailure, err) != nil {
		return
	}

	claims, _ := formats.AuthFromContext(r.Context())

	var uid int64 = 0
	if claims != nil {
		uid = claims.Uid
	}

	room, reconnect := game_logic.LoadOrStoreSinglePlayRoom(
		game_logic.NewSinglePlayerRoom(r, uid, conn))

	LogInfo(0, "WS: connected", r)
	go room.Run(r, reconnect)
}

func (h *SinglePlayHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
