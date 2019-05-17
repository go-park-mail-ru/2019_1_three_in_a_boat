package game_logic

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
)

type MultiPlayerRoom struct {
	RoomId     RoomId
	Conn1      atomicConn
	Conn2      atomicConn
	Uid1       int64
	Uid2       int64
	Score1     int64
	Score2     int64
	LastInput1 *Input
	LastInput2 *Input
	Request1   *http.Request
	Request2   *http.Request
	Snapshot   Snapshot
	Running    bool
}

func NewMultiPlayerRoom(
	r *http.Request, uid int64, conn *websocket.Conn) *MultiPlayerRoom {
	room := &MultiPlayerRoom{
		Conn1:      NewAtomicConn(conn),
		Conn2:      NewAtomicConn(nil),
		Uid1:       uid,
		Uid2:       0,
		LastInput1: NewInput(-math.Pi / 2),
		LastInput2: NewInput(math.Pi / 2),
		Request1:   r,
		Request2:   nil,
		RoomId:     uuid.New().String(),
		Snapshot:   NewSnapshot(),
	}
	return room
}

func (mpr *MultiPlayerRoom) ConnectSecondPlayer(
	r *http.Request, uid int64, conn *websocket.Conn) {
	mpr.Request2 = r
	mpr.Uid2 = uid
	mpr.Conn2 = NewAtomicConn(conn)
}

func (mpr *MultiPlayerRoom) Id() RoomId {
	return mpr.RoomId
}

func (mpr *MultiPlayerRoom) Disconnect(_conn *atomicConn) {
	conn := _conn.Get()
	if conn != nil {
		_ = conn.Close()
	}
	_conn.Reset()
}

func (mpr *MultiPlayerRoom) Disconnect1() {
	mpr.Disconnect(&mpr.Conn1)
}

func (mpr *MultiPlayerRoom) Disconnect2() {
	mpr.Disconnect(&mpr.Conn2)
}

func (mpr *MultiPlayerRoom) connWriteJSON(
	v interface{}, _conn *atomicConn, r *http.Request) bool {
	conn := _conn.Get()
	if conn != nil {
		err := conn.WriteJSON(v)
		if handlers.WSHandleErrForward(
			r, formats.ErrWebSocketFailure, mpr.RoomId, err) != nil {
			mpr.Disconnect(_conn)
			return false
		}
		return true
	}

	return false
}

func (mpr *MultiPlayerRoom) connWriteText(
	message []byte, _conn *atomicConn, r *http.Request) bool {
	conn := _conn.Get()
	if conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if handlers.WSHandleErrForward(
			r, formats.ErrWebSocketFailure, mpr.RoomId, err) != nil {
			mpr.Disconnect(_conn)
			return false
		}
		return true
	}

	return false
}

// Does not handle invalid json well - treats it like a disconnect. So the
// assumption is, JSON always marshals. The logs will show if it's a disconnect
// or a JSON marshaling error.
func (mpr *MultiPlayerRoom) WriteJSON(
	v1 interface{}, v2 interface{}) (ok1 bool, ok2 bool) {
	return mpr.connWriteJSON(v1, &mpr.Conn1, mpr.Request1),
		mpr.connWriteJSON(v2, &mpr.Conn2, mpr.Request2)
}

func (mpr *MultiPlayerRoom) WriteSameJSON(
	v interface{}) (ok1 bool, ok2 bool) {
	return mpr.WriteJSON(v, v)
}

// Same as WriteJSON - invalid JSON = disconnect.
func (mpr *MultiPlayerRoom) connReadJSON(
	v interface{}, _conn *atomicConn, r *http.Request) bool {
	conn := _conn.Get()
	if conn != nil {
		err := conn.ReadJSON(v)
		if handlers.WSHandleErrForward(
			r, formats.ErrWebSocketFailure, mpr.RoomId, err) != nil {
			mpr.Disconnect(_conn)
			return false
		}
		return true
	}

	return false
}

func (mpr *MultiPlayerRoom) ReadJSON(
	v1 interface{}, v2 interface{}) (ok1 bool, ok2 bool) {
	return mpr.connReadJSON(v1, &mpr.Conn1, mpr.Request1),
		mpr.connReadJSON(v2, &mpr.Conn2, mpr.Request2)
}

func (mpr *MultiPlayerRoom) ReadInput() (ok1 bool, ok2 bool) {
	return mpr.ReadJSON(&mpr.LastInput1, &mpr.LastInput2)
}

func (mpr *MultiPlayerRoom) Tick() {
	mpr.WriteJSON(
		MultiPlayerSnapshotData{
			Over1:             mpr.Snapshot.State == StateOverPlayer1,
			Over2:             mpr.Snapshot.State == StateOverPlayer2,
			OtherAngle:        -1 * mpr.LastInput2.Angle(),
			Score1:            mpr.Score1,
			Score2:            mpr.Score2,
			Hexagons:          mpr.Snapshot.Hexagons,
			CursorCircleAngle: mpr.Snapshot.CursorCircleAngle,
		},
		MultiPlayerSnapshotData{
			Over1:             mpr.Snapshot.State == StateOverPlayer1,
			Over2:             mpr.Snapshot.State == StateOverPlayer2,
			OtherAngle:        -1 * mpr.LastInput1.Angle(),
			Score1:            mpr.Score1,
			Score2:            mpr.Score2,
			Hexagons:          mpr.Snapshot.Hexagons,
			CursorCircleAngle: mpr.Snapshot.CursorCircleAngle,
		},
	)

	isOver := mpr.Snapshot.IsMultiplayerGameOver(mpr.LastInput1, mpr.LastInput2)

	if isOver {
		mpr.FinishGame()
	} else {
		score := mpr.Snapshot.Update()

		if mpr.Snapshot.State != StateOverPlayer1 {
			mpr.Score1 = score
		}
		if mpr.Snapshot.State != StateOverPlayer2 {
			mpr.Score2 = score
		}
	}
}

func (mpr *MultiPlayerRoom) FinishGame() {
	mpr.WriteSameJSON(
		MultiPlayerSnapshotData{
			Over1:  true,
			Over2:  true,
			Score1: mpr.Score1,
			Score2: mpr.Score2,
		})
	mpr.Disconnect1()
	mpr.Disconnect2()
	handlers.WSLogInfo(mpr.Request1, "Closing socket", mpr.RoomId)

	err := db.UpdateScoreById(settings.DB(), mpr.Uid1, mpr.Score1)
	if err != nil {
		handlers.WSLogError(mpr.Request1, "Failed to write game result", mpr.RoomId, err)
	}
	err = db.UpdateScoreById(settings.DB(), mpr.Uid2, mpr.Score2)
	if err != nil {
		handlers.WSLogError(mpr.Request2, "Failed to write game result", mpr.RoomId, err)
	}

	// Game.Rooms.Delete(string(mpr.RoomId)) // these are only stored in the queue
}

func (mpr *MultiPlayerRoom) ReadLoop1() {
	for mpr.connReadJSON(mpr.LastInput1, &mpr.Conn1, mpr.Request1) {
		mpr.Conn2 = mpr.Conn2
	}
}

func (mpr *MultiPlayerRoom) ReadLoop2() {
	for mpr.connReadJSON(mpr.LastInput2, &mpr.Conn2, mpr.Request2) {
		mpr.Conn1 = mpr.Conn1
	}
}

func (mpr *MultiPlayerRoom) Run() {
	ok1 := mpr.connWriteText([]byte(string(mpr.RoomId)+" 1"), &mpr.Conn1, mpr.Request1)
	ok2 := mpr.connWriteText([]byte(string(mpr.RoomId)+" 2"), &mpr.Conn2, mpr.Request2)

	if !ok1 && !ok2 {
		return
	} // if at least one of them is ok, let them play.. why not?

	time.Sleep(time.Second * 3)

	mpr.Snapshot.State = StateRunning
	go mpr.ReadLoop1()
	go mpr.ReadLoop2()
	tick := time.Tick(Settings.TickDuration)
	for mpr.Snapshot.State != StateOverBoth {
		<-tick
		mpr.Tick()
	}
}

type MultiPlayerSnapshotData struct {
	Over1             bool      `json:"over1,omitempty"`
	Over2             bool      `json:"over2,omitempty"`
	OtherAngle        float64   `json:"otherAngle"`
	Score1            int64     `json:"score1"`
	Score2            int64     `json:"score2"`
	Hexagons          []Hexagon `json:"hexes"`
	CursorCircleAngle float64   `json:"cursorCircleAngle"`
}
