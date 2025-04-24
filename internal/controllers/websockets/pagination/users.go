package pagination

import (
	"Ecadr/internal/repository"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Request struct {
	Search         string     `json:"search"`
	AfterCreatedAt *time.Time `json:"after_created_at"`
	AfterID        *uint      `json:"after_id"`
	Limit          int        `json:"limit"`
}

func UsersWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		var req Request
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[pagination.UsersWebSocket] Ошибка при чтении: %v", err)
			continue
		}

		err = json.Unmarshal(msg, &req)
		if err != nil {
			log.Printf("[pagination.UsersWebSocket] Error while unmarshalling json: %v", err)
			continue
		}

		if req.Limit <= 0 || req.Limit > 100 {
			req.Limit = 20
		}

		users, err := repository.GetUsersWithPagination(req.Search, req.AfterCreatedAt, req.AfterID, req.Limit)
		if err != nil {
			conn.WriteJSON(gin.H{"error": "could not fetch users"})
			continue
		}

		conn.WriteJSON(gin.H{"users": users})
	}
}
