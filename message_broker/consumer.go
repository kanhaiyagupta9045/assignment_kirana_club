package message_broker

import (
	"encoding/json"
	"log"

	"github.com/kanhaiyagupta9045/kirana_club/internals/process"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

var (
	err error
)

func NewConsumer(rabbitMQURL, queueName string) (*Consumer, error) {
	once.Do(func() {

		conn, err = amqp.Dial(rabbitMQURL)
		if err != nil {
			return
		}

		ch, err = conn.Channel()
		if err != nil {
			return
		}

		_, err = ch.QueueDeclare(
			queueName, // queue name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var data Data
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("Failed to unmarshal job: %v", err)
				continue
			}

			log.Printf("Received job: %v", data)
			process.ProcessJob(data.JobId, data.Store_Visit)
		}
	}()

	log.Println("Consumer started. Waiting for messages...")
	select {}
}
