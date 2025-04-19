package brokers

import (
	"github.com/streadway/amqp"
	"log"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func ConnectToRabbitMq() error {
	var err error
	RabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
