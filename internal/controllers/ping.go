package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Ping(c *gin.Context) {
	fmt.Println(c)
	log.Println("PPPPOOOOOOOOONNNNNGGGGG")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
