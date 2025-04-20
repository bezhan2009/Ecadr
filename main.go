package main

import (
	"Ecadr/configs"
	"Ecadr/internal/routes"
	"Ecadr/internal/security"
	"Ecadr/internal/server"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title Ecadr API
// @version 1.0.0

// @description API Server for Ecadr Application
// @host localhost:6565
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("example.env")
		if err != nil {
			panic(errors.New(fmt.Sprintf("error loading .env file. Error is %s", err)))
		}
	}

	security.AppSettings, err = configs.ReadSettings()
	if err != nil {
		panic(err)
	}
	security.SetConnDB(security.AppSettings)

	err = logger.Init()
	if err != nil {
		panic(err)
	}

	err = db.ConnectToDB()
	if err != nil {
		panic(err)
	}

	err = db.Migrate()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	mainServer := new(server.Server)
	go func() {
		if err = mainServer.Run(security.AppSettings.AppParams.PortRun, routes.InitRoutes(router)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while starting HTTP Service: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Printf("\n%s\n", yellow("Start of service termination"))

	// Закрытие соединения с БД
	err = db.CloseDBConn()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error closing database connection: %s", err.Error()))
	}

	//err = db2.CloseUserDBConn()
	//if err != nil {
	//	fmt.Println(fmt.Sprintf("Error closing user database connection: %s", err.Error()))
	//}

	// Корректное завершение HTTP-сервера
	if err = mainServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error while termination HTTP Service: %s", err)
	} else {
		fmt.Println(green("HTTP-service termination successfully"))
	}

	fmt.Println(red("End of program completion"))
}
