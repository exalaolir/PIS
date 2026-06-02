package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsURL      = "ws://localhost:3000/ws"
	maxBackoff = 30 * time.Second
)

func dialWithBackoff() (*websocket.Conn, chan string) {
	delay := time.Second
	for attempt := 1; ; attempt++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err == nil {
			log.Printf("соединение установлено (попытка %d)", attempt)
			replies := startReadLoop(conn)
			return conn, replies
		}
		log.Printf(" попытка %d: %v, повтор через %v", attempt, err, delay)
		time.Sleep(delay)
		delay *= 2
		if delay > maxBackoff {
			delay = maxBackoff
		}
	}
}

func startReadLoop(conn *websocket.Conn) chan string {
	replies := make(chan string, 8)

	conn.SetPingHandler(func(message string) error {
		log.Println("ping от сервера -> pong")
		return conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(10*time.Second))
	})

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			replies <- string(msg)
		}
	}()

	return replies
}

func main() {
	conn, replies := dialWithBackoff()

	for i := 1; i <= 5; i++ {
		msg := []byte("клиент: " + strconv.Itoa(i))

		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("записано:", err, "— переподключение...")
			conn.Close()
			conn, replies = dialWithBackoff()
			i--
			continue
		}
		log.Println("ОТПРАВЛЕНО:", string(msg))

		select {
		case reply := <-replies:
			log.Println("ПОЛУЧЕНО:", reply)
		case <-time.After(5 * time.Second):
			log.Println("таймаут — переподключение...")
			conn.Close()
			conn, replies = dialWithBackoff()
			i--
			continue
		}

		time.Sleep(1 * time.Second)
	}

	err := conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "конец"))
	if err != nil {
		log.Println("закрытие:", err)
		conn.Close()
	}
}
