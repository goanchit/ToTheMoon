package config

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var QueueConnection *amqp.Connection

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ConnectToQueue() {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:5672/", os.Getenv("MQ_ADMIN"), os.Getenv("MQ_PASSWORD"), os.Getenv("MQ_HOST"))
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to Connect To RabbitMQ")

	QueueConnection = conn
}
