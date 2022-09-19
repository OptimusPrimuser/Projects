package rabbit_mq

import (
	"backend/handler"
	"context"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQ struct {
	Connetion *amqp.Connection
	Channel   *amqp.Channel
	Queue     *amqp.Queue
	QName     string
}

func (rmq *RMQ) Init(queueName string) {
	rmq.QName = queueName
	var err error
	rmq.Connetion, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rmq.Channel, err = rmq.Connetion.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	q, err := rmq.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rmq.Queue = &q
}

func (rmq *RMQ) Consume() error {
	delivery, err := rmq.Channel.Consume(
		rmq.QName,
		"defaultConsumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for d := range delivery {
		// var entry JobEntry
		// json.Unmarshal(d.Body, &entry)
		// fmt.Println(entry)
		// go globals.Handler.
		// handler.Handler{}.Handle(d)
		var detail handler.JobEntry
		json.Unmarshal(d.Body, &detail)
		fmt.Println(detail)
		fmt.Println(handler.HandleRequest(d.Body))
		d.Ack(true)
	}
	return nil
}

func (rmq *RMQ) Publish(message interface{}) error {
	messageByte, err := json.Marshal(message)
	if err != nil {
		return err
	}
	pubMessage := amqp.Publishing{
		ContentType: "text/plain",
		Body:        messageByte,
	}
	err = rmq.Channel.PublishWithContext(
		context.Background(),
		"",
		rmq.QName,
		false,
		false,
		pubMessage,
	)
	if err != nil {
		return err
	}
	return nil
}

func (rmq *RMQ) Close() {
	rmq.Channel.Close()
	rmq.Connetion.Close()
}
