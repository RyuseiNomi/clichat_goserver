package main

import (
	"log"
	"net/http"
)

var serverPort = ":8000"

func main() {

	// ルームの生成
	r := NewRoom()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello`))
	})
	http.Handle("/room", r)

	// ルームの開始。これ以降、ルームはコマンドが終了するまで入退場を監視する
	go r.Run()

	log.Println("Webサーバの起動")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
