package main

import (
	"log"
	"net/http"

	"github.com/YuNaga224/websocketChat/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	//forwardは他のクライアントに送信するメッセージを保持するチャネル
	forward chan *message
	//チャットルームに参加しようとするクライアントのためのチャネル
	join chan *client
	//チャットルームを退室するクライアントのためのチャネル
	leave chan *client
	//在室するすべてのクライアントを保持する
	clients map[*client]bool
	//チャットルームの操作ログを受け取る
	tracer trace.Tracer
	//avatarはアバター情報を取得する
	avatar Avatar
}

// 新規ルームの生成を行う関数
func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加")
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信", msg.Message)
			//参加中のすべてのクライアントにメッセージを送信
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
					r.tracer.Trace("メッセージを送信")
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("送信に失敗。クライアントをクリーンアップ")

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

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	//webページが閉じられた際にクライアントを削除
	defer func() { r.leave <- client }()
	// メッセージの送信処理を別のgoroutineで実行
	go client.write()
	client.read()
}
