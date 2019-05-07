package game_logic

import (
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

type Room interface {
	Run()
	Id() RoomId
}

type RoomId = string

type SinglePlayerRoom struct {
	Conn      atomicConn
	Uid       int64
	Score     int64
	RoomId    RoomId
	Snapshot  Snapshot
	LastInput *Input
	Request   *http.Request
	Running   bool
}

func NewSinglePlayerRoom(
	r *http.Request, p1uid int64, conn *websocket.Conn) *SinglePlayerRoom {
	room := &SinglePlayerRoom{
		Conn:      NewAtomicConn(conn),
		Uid:       p1uid,
		RoomId:    uuid.New().String(),
		Snapshot:  NewSnapshot(),
		LastInput: NewInput(-math.Pi / 2),
		Request:   r,
	}
	return room
}

func (spr *SinglePlayerRoom) Id() RoomId {
	return spr.RoomId
}

func (spr *SinglePlayerRoom) Disconnect() {
	conn := spr.Conn.Get()
	if conn != nil {
		_ = conn.Close()
	}
	spr.Conn.Reset()
}

func (spr *SinglePlayerRoom) Reconnect(conn *websocket.Conn) {
	prev := spr.Conn.Get()
	if prev != nil {
		panic("an attempt to reconnect a connected game")
	}
	spr.Conn.Load(conn)
	spr.Snapshot.State = StateRunning
}

// Does not handle invalid json well - treats it like a disconnect. So the
// assumption is, JSON always marshals. The logs will show if it's a disconnect
// or a JSON marshaling error.
func (spr *SinglePlayerRoom) WriteJSON(v interface{}) bool {
	conn := spr.Conn.Get()
	if conn != nil {
		err := conn.WriteJSON(v)
		if handlers.WSHandleErrForward(
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
	conn := spr.Conn.Get()
	if conn != nil {
		err := conn.ReadJSON(v)
		if handlers.WSHandleErrForward(
			spr.Request, formats.ErrWebSocketFailure, spr.RoomId, err) != nil {
			spr.Disconnect()
			return false
		}
		return true
	}

	return false
}

func (spr *SinglePlayerRoom) ReadInput() bool {
	return spr.ReadJSON(&spr.LastInput)
}

// Errors are mostly logged and ignored - the game goes on as if nothing
// happened, anticipating a reconnect, since all errors boil down to network
// or json errors, the latter being a programming error, which is properly
// logged but otherwise treated like a disconnect. Gorilla WS implementation
// doesn't really provide a way to differ between the two other than the
// Error() method, so that's why all errors are treated as disconnects.
func (spr *SinglePlayerRoom) Tick() {
	spr.WriteJSON(
		SinglePlayerSnapshotData{
			Over:              false,
			Score:             spr.Score,
			Hexagons:          spr.Snapshot.Hexagons,
			CursorCircleAngle: spr.Snapshot.CursorCircleAngle,
		})

	isOver := spr.Snapshot.IsGameOver(spr.LastInput)

	if isOver {
		spr.FinishGame()
	} else {
		spr.Score = spr.Snapshot.Update()
	}
}

func (spr *SinglePlayerRoom) FinishGame() {
	spr.WriteJSON(
		SinglePlayerSnapshotData{Over: true, Score: spr.Score})
	spr.Disconnect()
	handlers.WSLogInfo(spr.Request, "Closing socket", spr.RoomId)
	err := db.UpdateScoreById(settings.DB(), spr.Uid, spr.Score)
	if err != nil {
		handlers.WSLogError(spr.Request, "Failed to write game result", spr.RoomId, err)
	}
	Game.Rooms.Delete(string(spr.RoomId))
}

func (spr *SinglePlayerRoom) ReadLoop() {
	for spr.ReadInput() {
	}
}

func (spr *SinglePlayerRoom) Run() {
	conn := spr.Conn.Get()
	if conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, []byte(spr.RoomId))
		if err != nil {
			handlers.LogError(0, "WS: unexpected disconnect", spr.Request)
			Game.Rooms.Delete(string(spr.RoomId))
			return
		}
	} else {
		Game.Rooms.Delete(string(spr.RoomId))
	}

	spr.Snapshot.State = StateRunning
	go spr.ReadLoop()
	tick := time.Tick(Settings.TickDuration)
	for spr.Snapshot.State != StateOverBoth {
		<-tick
		spr.Tick()
	}
}

type SinglePlayerSnapshotData struct {
	Over              bool      `json:"over,omitempty"` // false is omitted
	Score             int64     `json:"score"`
	Hexagons          []Hexagon `json:"hexes"`
	CursorCircleAngle float64   `json:"cursorCircleAngle"`
}
