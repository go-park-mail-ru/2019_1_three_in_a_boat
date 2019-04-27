package chat_logic

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/chat/messages"
	"sync"
)

type Chat struct {
	Sockets sync.Map
}

var MainChat = Chat{}

func WriteToAll(message messages.Message) {
	MainChat.Sockets.Range(func(_sockId, _sock interface{}) bool {
		sock := _sock.(*ChatSocket)
		sock.WriteJSON(message)
		return true
	})
}

func GetLastMessages() {

}
