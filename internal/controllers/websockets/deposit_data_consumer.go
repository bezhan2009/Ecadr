package controllers

import (
	"Deposit/pkg/brokers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // разрешаем любые источники (на проде лучше ограничить)
	},
}

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Ошибка при апгрейде соединения:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}
	}
}

func StartQueueConsumer() {
	msgs, err := brokers.RabbitChannel.Consume(
		"deposit_portal", // очередь
		"",               // consumer tag
		false,            // autoAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // args
	)
	if err != nil {
		log.Fatal("Ошибка при получении сообщений:", err)
	}

	go func() {
		for msg := range msgs {
			log.Println("Получено сообщение:", string(msg.Body))

			mu.Lock()
			for conn := range clients {
				err := conn.WriteMessage(websocket.TextMessage, msg.Body)
				if err != nil {
					log.Println("Ошибка при отправке клиенту:", err)
					conn.Close()
					delete(clients, conn)
				}
			}
			mu.Unlock()

			// Не подтверждаем сообщение
		}
	}()
}
