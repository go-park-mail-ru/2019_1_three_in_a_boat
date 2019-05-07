package routes

import (
	"fmt"
	"net/http"

	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/game/game_logic"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
)

type MultiPlayHandler struct{}

func (h *MultiPlayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		LogError(1, fmt.Sprintf(
			"Connection %s: (%s)", formats.ErrWebSocketFailure, err.Error()), r)
		return
	}

	claims, _ := formats.AuthFromContext(r.Context())

	var uid int64 = 0
	if claims != nil {
		uid = claims.Uid
	}

	room, ready := game_logic.GetOrCreateMPRoom(r, uid, conn)
	if ready {
		room.ConnectSecondPlayer(r, uid, conn)
		room.Run()
	} // else do nothing - when another player connects, it will be handled there
}

func (h *MultiPlayHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
