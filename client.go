package main

import (
	"flag"
	"log"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "go-chat-server:8080", "http service address")

// Client チャットを行う1人のユーザを表す
type Client struct {
	name string // クライアントの名前を格納するフィールド
	// クライアントのためのWebSocket
	socket *websocket.Conn
	// メッセージが送られるチャネル
	send chan []byte
	// クライアントが参加するチャットルーム
	room *room
}

// Read クライアントのWebsocketからデータの読み込みを行う
func (c *Client) Read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// Write 継続的にsendチャネルからメッセージを受け取り、Websocketへの書き込みを行う
func (c *Client) Write() {
	for msg := range c.send {
		log.Println("Websocketへの書き込みを実行")
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			// TODO クライアントへの転送処理を書く
			break
		}
	}
	c.socket.Close()
}
