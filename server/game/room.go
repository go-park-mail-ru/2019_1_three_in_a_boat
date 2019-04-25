package game

import (
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings"
)

type Room interface {
	Run(*http.Request, bool)
	Id() RoomId
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
	LastInput *Input
	Request   *http.Request
	Running   bool
}

func (spr *SinglePlayerRoom) Id() RoomId {
	return spr.RoomId
}

func (spr *SinglePlayerRoom) Disconnect() {
	spr.Conn = nil
	if spr.Snapshot.State == StateRunning {
		spr.Snapshot.State = StateDisconnect
	}
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
	if spr.ReadJSON(&spr.LastInput) {
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
		LastInput: NewInput(-math.Pi / 2),
		Request:   r,
	}
	return room
}

// Errors are mostly logged and ignored - the game goes on as if nothing
// happened, anticipating a reconnect, since all errors boil down to network
// or json errors, the latter being a programming error, which is properly
// logged but otherwise treated like a disconnect. Gorilla WS implementation
// doesn't really provide a way to differ between the two other than the
// Error() method, so that's why all errors are treated as disconnects.
func (spr *SinglePlayerRoom) Tick() (isOver bool) {
	spr.Snapshot.State = StateRunning

	spr.WriteJSON(
		SinglePlayerSnapshotData{
			Over:              false,
			Score:             spr.Snapshot.Score,
			Hexagons:          spr.Snapshot.Hexagons,
			CursorCircleAngle: spr.Snapshot.CursorCircleAngle,
		})

	isOver = spr.Snapshot.Update(spr.LastInput)

	if isOver {
		spr.FinishGame()
	}

	return
}

func (spr *SinglePlayerRoom) FinishGame() {
	spr.WriteJSON(
		SinglePlayerSnapshotData{Over: true, Score: spr.Snapshot.Score})
	if spr.Conn != nil {
		_ = spr.Conn.Close()
	}
	spr.Disconnect()
	WSLogInfo(spr.Request, "Closing socket", spr.RoomId)
	err := db.UpdateScoreById(settings.DB(), spr.Uid, spr.Snapshot.Score)
	if err != nil {
		WSLogError(spr.Request, "Failed to write game result", spr.RoomId, err)
	}
	Game.Rooms.Delete(string(spr.RoomId))
}

func (spr *SinglePlayerRoom) ReadLoop() {
	for spr.ReadInput() {
	}
}

func (spr *SinglePlayerRoom) Run(r *http.Request, reconnect bool) {
	go spr.ReadLoop()
	if spr.Connected() {
		err := spr.Conn.WriteMessage(websocket.TextMessage, []byte(spr.RoomId))
		if err != nil {
			LogError(0, "WS: unexpected disconnect", r)
		}
	}

	tick := time.Tick(Settings.TickDuration)
	for spr.Snapshot.State != StateOver {
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
