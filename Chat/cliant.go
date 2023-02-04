package main

import (
	"time"

	"github.com/gorilla/websocket"
)

//cliantはチャットを行うユーザーを表す

type client struct {
	//websocketへの参照
	socket *websocket.Conn
	//受信したメッセージが待ち行列のように蓄積されるバッファ付きのチャネル
	send chan *message
	room *room
	//userDataはユーザーに関する情報を保持する
	userData map[string]interface{}
}

// サーバー皮の処理
func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// クライアント側の処理
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
