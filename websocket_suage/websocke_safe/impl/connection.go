package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type WsConnection struct {
	conn      *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	isClosed  bool
	closeLock sync.Mutex
}

func InitConnection(conn *websocket.Conn) (ws *WsConnection, err error) {
	ws = &WsConnection{
		conn:      conn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	// 拉起读与写的协程
	go ws.readLoop()
	go ws.writeLoop()
	return
}

// 如果已经关闭，有可能会卡死在这里，所以这里也要进行通知
func (conn *WsConnection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")

	}
	return
}

func (conn *WsConnection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")

	}
	return
}

func (conn *WsConnection) Close() {
	// Close()本身是线程安全的，并且可重入的。
	conn.Close()
	// 如果多次close会出问题，这时候需要一个flag 同时需要一个锁来保护它
	conn.closeLock.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.closeLock.Unlock()
}

func (conn *WsConnection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.conn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR

		}
	}
ERR:
	conn.Close()
}

// 如果写失败关闭了连接，但是如果读的队列满了，这时候readLoop协程会一直卡在哪儿，造成内存泄漏了。所以需要一个通知消息
func (conn *WsConnection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR

		}
		if err = conn.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
