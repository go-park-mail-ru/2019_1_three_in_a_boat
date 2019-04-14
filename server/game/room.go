package game

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
)

type Room interface {
	Run(*http.Request, bool)
}

type RoomId = string

type MultiPlayerRoom struct {
	Player1 *SinglePlayerRoom
	Player2 *SinglePlayerRoom
}

func (r *MultiPlayerRoom) IsSinglePlayer() bool {
	return r.Player2 == nil
}

type SinglePlayerRoom struct {
	Conn      *websocket.Conn
	Uid       int64
	RoomId    RoomId
	Snapshot  Snapshot
	LastInput Input
	Request   *http.Request
	Running   bool
}

func (spr *SinglePlayerRoom) Disconnect() {
	spr.Conn = nil
	spr.Snapshot.State = StateDisconnect
}

func (spr *SinglePlayerRoom) Connected() bool {
	return spr.Conn != nil
}

func (spr *SinglePlayerRoom) Reconnect(conn *websocket.Conn) {
	if spr.Connected() {
		panic("an attempt to reconnect a connected game")
	}
	spr.Conn = conn
	spr.Snapshot.State = StateRunning
}

// Does not handle invalid json well - treats it like a disconnect. So the
// assumption is, JSON always marshals. The logs will show if it's a disconnect
// or a JSON marshaling error.
func (spr *SinglePlayerRoom) WriteJSON(v interface{}) bool {
	if spr.Connected() {
		err := spr.Conn.WriteJSON(v)
		if WSHandleErrForward(
			spr.Request, formats.ErrWebSocketFailure, spr.RoomId, err) != nil {
			spr.Disconnect()
			return false
		} else {
			return true
		}
	}

	return false
}

// Same as WriteJSON - invalid JSON = disconnect.
func (spr *SinglePlayerRoom) ReadJSON(v interface{}) bool {
	if spr.Connected() {
		err := spr.Conn.ReadJSON(v)
		if WSHandleErrForward(
			spr.Request, formats.ErrWebSocketFailure, spr.RoomId, err) != nil {
			spr.Disconnect()
			return false
		} else {
			return true
		}
	}

	return false
}

func (spr *SinglePlayerRoom) ReadInput() bool {
	input := Input{}
	if spr.ReadJSON(&input) {
		spr.LastInput = input
		return true
	}

	return false
}

func NewSinglePlayerRoom(
	r *http.Request, p1uid int64, conn *websocket.Conn) *SinglePlayerRoom {
	room := &SinglePlayerRoom{
		Conn:      conn,
		Uid:       p1uid,
		RoomId:    uuid.New().String(),
		Snapshot:  NewSnapshot(),
		LastInput: Input{math.Pi / 2},
		Request:   r,
	}
	Game.Rooms.Store(room.RoomId, room)
	return room
}

// Errors are mostly logged and ignored - the game goes on as if nothing
// happened, anticipating a reconnect, since all errors boil down to network
// or json errors, the latter being a programming error, which is properly
// logged but otherwise treated like a disconnect. Gorilla WS implementation
// doesn't really provide a way to differ between the two other than the
// Error() method, so that's why all errors are treated as disconnects.
func (spr *SinglePlayerRoom) Tick() (
	isOver bool) {
	if !spr.ReadInput() {
		WSLogError(spr.Request, "WS: unexpected disconnect",
			spr.RoomId, errors.New(formats.ErrWebSocketFailure))
	}
	spr.Snapshot.State = StateRunning
	if isOver {
		spr.WriteJSON(
			SinglePlayerSnapshotData{Over: true, Score: spr.Snapshot.Score})
		return
	} else {
		spr.WriteJSON(
			SinglePlayerSnapshotData{
				Over:     false,
				Score:    spr.Snapshot.Score,
				Hexagons: spr.Snapshot.Hexagons,
			})
	}

	isOver = spr.Snapshot.Update(spr.LastInput)

	if isOver {
		spr.WriteJSON(
			SinglePlayerSnapshotData{Over: true, Score: spr.Snapshot.Score})

	}

	return
}

func (spr *SinglePlayerRoom) Run(r *http.Request, reconnect bool) {
	spr.LastInput = Input{Angle: math.Pi / 2}
	err := spr.Conn.WriteMessage(websocket.TextMessage, []byte(spr.RoomId))
	if err != nil {
		LogError(0, "WS: unexpected disconnect", r)
	}

	tick := time.Tick(Settings.TickDuration)
	for spr.Snapshot.State != StateOver {
		<-tick
		spr.Tick()
	}
}

type SinglePlayerSnapshotData struct {
	Over     bool      `json:"over,omitempty"` // false is omitted
	Score    int64     `json:"score"`
	Hexagons []Hexagon `json:"hexes"`
}
