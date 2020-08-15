package main

import (
	"flag"

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

// Read WebSocketへの書き込みを行う
// ここでは、全てのクライアントの送信済みメッセージを読み込んでいる
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

// Write WebSocketへの書き込みを行う。
// ここでは主に、メッセージを送信する際にルームに在籍しているメンバー全員に転送をする処理を読んでいる
func (c *Client) Write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
