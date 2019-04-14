package routes

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// A handler that handles a ~multiple~ users resource. The handler itself is
// simply a struct that forwards the requests, depending on the method to one of
// GetUsers or PostUsers. Accepts PUT and POST requests. Implements
// route.Handler. For details on what do both resources do, see
// users_(method).go.
type PlayHandler struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *PlayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
}

func (h *PlayHandler) Settings() map[string]RouteSettings {
	return map[string]RouteSettings{
		"POST": {
			AuthRequired:           true,
			CorsAllowed:            true,
			CsrfProtectionRequired: true,
		},
	}
}
