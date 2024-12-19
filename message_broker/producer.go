package message_broker

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/kanhaiyagupta9045/kirana_club/internals/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Data struct {
	JobId       primitive.ObjectID
	Store_Visit models.StoresVisit
}
type Producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

var (
	once    sync.Once
	conn    *amqp.Connection
	initErr error
	ch      *amqp.Channel
)

func NewProducer(rabbitMQURL, queueName string) (*Producer, error) {
	once.Do(func() {
		conn, initErr = amqp.Dial(rabbitMQURL)
		if initErr != nil {
			return
		}

		ch, initErr = conn.Channel()
		if initErr != nil {
			return
		}

		_, initErr = ch.QueueDeclare(
			queueName, // queue name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if initErr != nil {
			return
		}

	})
	if initErr != nil {
		return nil, initErr
	}
	return &Producer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (p *Producer) Publish(data Data) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published job: %v", data)
	return nil
}
