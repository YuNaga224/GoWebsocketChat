package main

import (
	"github.com/gorilla/websocket"
)

//cliantはチャットを行うユーザーを表す

type client struct {
	//websocketへの参照
	socket *websocket.Conn
	//受信したメッセージが待ち行列のように蓄積されるバッファ付きのチャネル
	send chan []byte
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
