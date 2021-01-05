package app

import (
	"fmt"
	"log"

	"github.com/hamba/avro"
	"github.com/marktsoy/hb_sender/internal/models"
	"github.com/streadway/amqp"
)

func ListenToSubscribers(c *Config, cb func(models.Subscription)) {

	conn, err := amqp.Dial(c.RabbitAddr)
	if err != nil {
		panic("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		c.RabbitSubscriberQueue, // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		panic("Failed to declare queue")
	}

	subscribers, err := ch.Consume(q.Name, "subscription-listener", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case s := <-subscribers:
			fmt.Println(s.Body)
			out := models.Subscription{}
			schema := models.SchemaSubscription()
			err = avro.Unmarshal(schema, s.Body, &out)
			if err != nil {
				log.Fatal("Wrong schema")
			}
			fmt.Printf("%v \n %v", out.ChatID, out.IsSubscribed)
			cb(out)
			fmt.Printf("Chat subscribed \t %v", out.ChatID)
			fmt.Printf("Received a message: %s\n", s.Body)
		}
	}
}
