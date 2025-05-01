// routes/workers.go

package routes

import (
	"Ecadr/internal/controllers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

type workerTask func()

var (
	workerPool   chan workerTask
	heavyLimiter *rate.Limiter
)

func init() {
	workerPool = make(chan workerTask, 100)
	for i := 0; i < 10; i++ {
		go func(id int) {
			log.Printf("[worker %d] started", id)
			for task := range workerPool {
				task()
				log.Printf("[worker %d] task done", id)
			}
		}(i)
	}

	heavyLimiter = rate.NewLimiter(5, 10) // 5 rps, burst 10
}

// TaskSync — handler, но возвращает (результат, ошибка)
type TaskSync func(*gin.Context) (interface{}, error)

func runHeavyTask(c *gin.Context, handler TaskSync) {
	// rate limiter
	if !heavyLimiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	// создаём канал для результата
	resultCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	// копируем контекст (для безопасного чтения)
	cCp := c.Copy()

	// ставим задачу в пул
	workerPool <- func() {
		res, err := handler(cCp)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- res
	}

	// ждём результата или таймаута
	select {
	case res := <-resultCh:
		c.JSON(http.StatusOK, res)
	case err := <-errCh:
		controllers.HandleError(c, err)
	case <-time.After(30 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout"})
	}
}
