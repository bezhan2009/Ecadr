package main

import (
	"Ecadr/configs"
	"Ecadr/internal/routes"
	"Ecadr/internal/security"
	"Ecadr/internal/server"
	"Ecadr/pkg/brokers"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Ecadr API
// @version 1.0.0

// @description API Server for Ecadr Application
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// цветной вывод
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// загрузка конфигов
	if err := godotenv.Load(".env"); err != nil {
		if err = godotenv.Load("example.env"); err != nil {
			log.Fatalf("error loading .env file: %v", err)
		}
	}

	// конфиги
	var err error
	security.AppSettings, err = configs.ReadSettings()
	if err != nil {
		log.Fatalf("failed to read settings: %v", err)
	}
	security.SetConnDB(security.AppSettings)

	// логгер
	if err = logger.Init(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	// подключение к БД (2 попытки)
	if err = db.ConnectToDB(); err != nil {
		time.Sleep(10 * time.Second)
		if err = db.ConnectToDB(); err != nil {
			log.Fatalf("failed to connect to DB: %v", err)
		}
	}
	if err = db.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Redis и RabbitMQ
	if err = db.InitializeRedis(security.AppSettings.RedisParams); err != nil {
		log.Fatalf("redis init failed: %v", err)
	}
	if err = brokers.ConnectToRabbitMq(security.AppSettings.RabbitParams); err != nil {
		log.Fatalf("rabbitmq connect failed: %v", err)
	}

	// роутинг
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	mainServer := new(server.Server)

	// запуск сервера в отдельной горутине
	certFile := os.Getenv("SSL_CERT_PATH")
	keyFile := os.Getenv("SSL_KEY_PATH")

	go func() {
		if err = mainServer.Run(security.AppSettings.AppParams.PortRun, routes.InitRoutes(router), certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while starting HTTPS Service: %s", err)
		}
	}()

	// ловим сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println(yellow("\nStart of service termination"))

	// Останавливаем всё с timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = mainServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %s", err)
	} else {
		fmt.Println(green("HTTP-service terminated successfully"))
	}

	if err = db.CloseDBConn(); err != nil {
		log.Printf("DB close error: %s", err)
	}
	if err = db.CloseRedisConnection(); err != nil {
		log.Printf("Redis close error: %s", err)
	}

	fmt.Println(red("End of program completion"))
}
