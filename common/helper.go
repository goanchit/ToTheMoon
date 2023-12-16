package common

import (
	"context"
	"log"
	"taskmanager/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func createChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	ch, err := connection.Channel()
	if err != nil {
		log.Fatalln(err.Error())
		return ch, err
	}
	return ch, nil
}

func PublishToQueue(qName string, message string) error {
	qConnection := config.QueueConnection

	ch, err := createChannel(qConnection)

	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		qName,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalln("Failed to publish a message")
		return err
	}
	log.Printf("Published message to Queue %s", qName)
	return nil
}

func QueueConsumer(qName string) {
	qConnection := config.QueueConnection

	ch, err := createChannel(qConnection)

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		qName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln(err.Error())
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Listening for messages")
	<-forever

}
