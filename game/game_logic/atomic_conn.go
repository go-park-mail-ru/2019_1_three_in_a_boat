package game_logic

import (
	"sync/atomic"
	"unsafe"

	"github.com/gorilla/websocket"
)

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
