package app

import (
	"fmt"
	"log"
	"time"

	"github.com/hamba/avro"

	"github.com/marktsoy/hb_sender/internal/models"
	"github.com/streadway/amqp"
)

// ListenToMessages on specific queueName and perfoms callback to the message
func ListenToMessages(c *Config, cb func(models.Message)) {

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
		c.RabbitMessageQueue,            // name
		false,                           // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		amqp.Table{"x-max-priority": 5}, // arguments
	)
	if err != nil {
		panic("Failed to declare queue")
	}

	msgs, err := ch.Consume(q.Name, "message-listener", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	timeoutDuration := time.Duration(c.RabbitMessageTimeout) * time.Second
	go func() {
		for {
			time.Sleep(timeoutDuration)
			select {
			case msg := <-msgs:
				fmt.Println(msg.Body)
				out := models.Message{}
				schema := models.SchemaMessage()
				err = avro.Unmarshal(schema, msg.Body, &out)
				if err != nil {
					log.Fatal("Wrong schema")
				}
				fmt.Printf("%v \n", out.Content)
				cb(out)
				fmt.Printf("Received a message: %s\n", msg.Body)
			}
		}
	}()

	<-forever
}
