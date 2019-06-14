package main

import (
	"github.com/emacsvi/dogolang/websocket_suage/websocke_safe/impl"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	// 允许所有CORS跨域请求
	upgrade = websocket.Upgrader{
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		conn   *impl.WsConnection
		data   []byte
	)

	if wsConn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		return
	}
	defer conn.Close()

	// 启动一个协程，一秒发送一次心跳
	go func() {
		var (
			err error
		)
		for {
			time.Sleep(time.Second)
			if err = conn.WriteMessage([]byte("heart beat")); err != nil {
				return
			}
		}
	}()

	// 不停读数据并且原样返回给客户端
	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandle)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
