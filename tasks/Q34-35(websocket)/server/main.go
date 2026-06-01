package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 20 * time.Second
	pingPeriod = 5 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("WebSocket-сервер: http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	log.Println("соединение установлено")

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		log.Println("pong от клиента")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	outgoing := make(chan []byte, 8)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		readLoop(conn, outgoing)
		wg.Done()
	}()

	go func() {
		writeLoop(conn, outgoing)
		wg.Done()
	}()

	wg.Wait()
	conn.Close()
}

func readLoop(conn *websocket.Conn, outgoing chan<- []byte) {
	log.Println("goroutine чтения запущена — только чтение")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("goroutine чтения: клиент закрыл соединение")
			} else {
				log.Println("goroutine чтения: ошибка чтения:", err)
			}
			close(outgoing)
			return
		}
		log.Println("goroutine чтения: ПОЛУЧЕНО:", string(msg))
		outgoing <- msg
	}
}

func writeLoop(conn *websocket.Conn, outgoing <-chan []byte) {
	log.Println("goroutine записи запущена — только запись")
	sendPing := func() bool {
		if err := conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(writeWait)); err != nil {
			log.Println("goroutine записи: ping:", err)
			return false
		}
		log.Println("goroutine записи: ping отправлен")
		return true
	}
	if !sendPing() {
		return
	}
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			if !sendPing() {
				ticker.Stop()
				return
			}
		case msg, ok := <-outgoing:
			if !ok {
				ticker.Stop()
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("goroutine записи: ошибка записи:", err)
				ticker.Stop()
				return
			}
			log.Println("goroutine записи: ОТПРАВЛЕНО:", string(msg))
		}
	}
}
