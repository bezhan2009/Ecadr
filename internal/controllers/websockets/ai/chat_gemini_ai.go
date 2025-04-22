package aiWebSocket

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiService "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/brokers"
	"Ecadr/pkg/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // TODO: restrict origins in production
		},
	}

	clients         = make(map[*websocket.Conn]bool)
	mu              sync.Mutex
	pendingMu       sync.Mutex
	pendingMessages = make(map[uint64]amqp.Delivery)
)

func StartQueueConsumer(userID uint) error {
	_, err := brokers.RabbitChannel.QueueDeclare(
		"ai_chat",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error.Printf("[StartQueueConsumer] Queue declare error: %v", err)
		return err
	}

	msgs, err := brokers.RabbitChannel.Consume(
		"ai_chat",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error.Printf("[StartQueueConsumer] Queue consume error: %v", err)
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handleMessage(userID, msg); err != nil {
				logger.Error.Printf("[handleMessage] error: %v", err)
			}
		}
	}()

	return nil
}

func handleMessage(userID uint, msg amqp.Delivery) error {
	log.Printf("Получено сообщение: %s", string(msg.Body))

	pendingMu.Lock()
	pendingMessages[msg.DeliveryTag] = msg
	pendingMu.Unlock()

	user, err := service.GetUserByID(userID)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	if err := service.CreateMessage(models.Message{
		Text:   string(msg.Body),
		UserID: userID,
	}); err != nil {
		log.Printf(err.Error())
		return err
	}

	respAI, err := aiService.SendMessageToGeminiAI(
		user.KindergartenNotes,
		user.SchoolGrades,
		user.Achievements,
		string(msg.Body),
	)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	messageWithID := map[string]interface{}{
		"tag":  msg.DeliveryTag,
		"body": respAI,
	}
	payload, err := json.Marshal(messageWithID)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	sendToAllClients(payload)
	return nil
}

func sendToAllClients(payload []byte) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("[sendToAllClients] Ошибка при отправке: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

func ChatAIWebSocketHandler(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ChatAIWebSocketHandler] Ошибка при апгрейде: %v", err)
		return
	}
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	if err := StartQueueConsumer(userID); err != nil {
		sendError(conn, "Internal Server Error")
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[ChatAIWebSocketHandler] Ошибка при чтении: %v", err)
			break
		}
		log.Printf("Получено от клиента: %s", string(msg))

		err = brokers.RabbitChannel.Publish(
			"",        // default exchange
			"ai_chat", // queue name
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			},
		)
		if err != nil {
			log.Printf("[ChatAIWebSocketHandler] Ошибка при публикации в очередь: %v", err)
		}
	}
}

func sendError(conn *websocket.Conn, message string) {
	errorPayload := map[string]string{"error": message}
	payload, _ := json.Marshal(errorPayload)
	conn.WriteMessage(websocket.TextMessage, payload)
}
