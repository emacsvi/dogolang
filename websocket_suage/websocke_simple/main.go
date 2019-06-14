package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		conn *websocket.Conn
		err  error
		data []byte
	)
	if conn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}
	defer conn.Close()

	for {
		if _, data, err = conn.ReadMessage(); err != nil {
			return
		}

		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandle)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
