package routes

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/game"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
)

// A handler that handles a ~multiple~ users resource. The handler itself is
// simply a struct that forwards the requests, depending on the method to one of
// GetUsers or PostUsers. Accepts PUT and POST requests. Implements
// route.Handler. For details on what do both resources do, see
// users_(method).go.
type SinglePlayHandler struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *SinglePlayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if HandleErrForward(w, r, formats.ErrWebSocketFailure, err) != nil {
		return
	}

	// claims, _ := formats.AuthFromContext(r.Context())

	room, reconnect := game.LoadOrStoreSinglePlayRoom(
		0, game.NewSinglePlayerRoom(r, 0, conn))

	LogInfo(0, "WS: connected", r)
	go room.Run(r, reconnect)
}

func (h *SinglePlayHandler) Settings() map[string]RouteSettings {
	return map[string]RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: true,
		},
	}
}