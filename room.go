package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte      //他のクライアントに転送するためのメッセージを保持するチャネル
	join    chan *Client     //ルームに参加しようとしているクライアントのためのチャネル
	leave   chan *Client     //ルームから退出しようとしているクライアントのためのチャネル
	clients map[*Client]bool //在室している全てのクライアントを保持
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// Run チャットルームの開始および全てのチャネルの監視を開始
func (r *room) Run() {
	log.Println("チャットルームを開始")
	for {
		select {
		case client := <-r.join:
			// ルームへの参加
			str := fmt.Sprintf("%sさん が参加しました", client.name)
			log.Println(str)
			r.clients[client] = true
			for client := range r.clients {
				client.send <- []byte(str)
			}
		case client := <-r.leave:
			// ルームからの退出
			str := fmt.Sprintf("%sさん が退出しました", client.name)
			log.Println(str)
			for client := range r.clients {
				client.send <- []byte(str)
			}
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			log.Println("メッセージを転送:", string(msg))
			// 全てのクライアントへのメッセージの転送
			for client := range r.clients {
				select {
				case client.send <- msg: //TODO 送信者名もmsgに含める
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

// ServeHTTP ルームのアップグレードを行う
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	// QueryStringの抽出
	queryString := req.URL.Query()
	value, _ := queryString["name"]

	if err != nil {
		log.Fatal("ServeHttp: ", err)
		return
	}
	client := &Client{
		name:   value[0],
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.Write()
	client.Read()
}

// NewRoom 初期値のルームを返却
func NewRoom() *room {
	log.Println("チャットルームをルームを作成します")
	return &room{
		forward: make(chan []byte),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}
