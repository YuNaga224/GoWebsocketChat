package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	//forwardは他のクライアントに送信するメッセージを保持するチャネル
	forward chan []byte
	//チャットルームに参加しようとするクライアントのためのチャネル
	join chan *client
	//チャットルームを退室するクライアントのためのチャネル
	leave chan *client
	//在室するすべてのクライアントを保持する
	clients map[*client]bool
}

// 新規ルームの生成を行う関数
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			//参加中のすべてのクライアントにメッセージを送信
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)

				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// WebSocketを使用するためにHTTP接続をアップグレード
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	//webページが閉じられた際にクライアントを削除
	defer func() { r.leave <- client }()
	// メッセージの送信処理を別のgoroutineで実行
	go client.write()
	client.read()
}
