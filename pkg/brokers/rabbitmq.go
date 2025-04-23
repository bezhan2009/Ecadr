package brokers

import (
	"Ecadr/internal/app/models"
	"fmt"
	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func ConnectToRabbitMq(params models.RabbitParams) error {
	var err error
	fmt.Println(params.URLConn)
	RabbitConn, err = amqp.Dial(params.URLConn)
	if err != nil {
		RabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			return err
		}
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		return err
	}

	return nil
}
