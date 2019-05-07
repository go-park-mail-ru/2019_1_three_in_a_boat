package game_logic

import (
	"math"
	"net/http"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
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

type atomicConn struct {
	p unsafe.Pointer
}

//noinspection GoExportedFuncWithUnexportedType
func NewAtomicConn(conn *websocket.Conn) atomicConn {
	ac := atomicConn{}
	ac.Load(conn)
	return ac
}

func (ac *atomicConn) Get() *websocket.Conn {
	return (*websocket.Conn)(atomic.LoadPointer(&ac.p))
}

func (ac *atomicConn) Reset() {
	atomic.StorePointer(&ac.p, nil)
}

func (ac *atomicConn) Load(conn *websocket.Conn) {
	atomic.StorePointer(&ac.p, unsafe.Pointer(conn))
}

type SinglePlayerRoom struct {
	Conn      atomicConn
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
	spr.Conn.Reset()
	if spr.Snapshot.State == StateRunning {
		spr.Snapshot.State = StateDisconnect
	}
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
	conn := spr.Conn.Get()
	if conn != nil {
		_ = conn.Close()
	}
	spr.Disconnect()
	handlers.WSLogInfo(spr.Request, "Closing socket", spr.RoomId)
	err := db.UpdateScoreById(settings.DB(), spr.Uid, spr.Snapshot.Score)
	if err != nil {
		handlers.WSLogError(spr.Request, "Failed to write game result", spr.RoomId, err)
	}
	Game.Rooms.Delete(string(spr.RoomId))
}

func (spr *SinglePlayerRoom) ReadLoop() {
	for spr.ReadInput() {
	}
}

func (spr *SinglePlayerRoom) Run(r *http.Request, reconnect bool) {
	go spr.ReadLoop()
	conn := spr.Conn.Get()
	if conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, []byte(spr.RoomId))
		if err != nil {
			handlers.LogError(0, "WS: unexpected disconnect", r)
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
