package app

import (
	"fmt"
	"log"
	"time"

	"github.com/hamba/avro"

	"github.com/marktsoy/hb_sender/internal/models"
	"github.com/streadway/amqp"
)

var (
	priorities = []string{"high", "medium", "low"}
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

	err = ch.ExchangeDeclare(
		"messages_topic", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		panic("Failed to declare exchange")
	}

	qs := make(map[string](<-chan amqp.Delivery))

	for _, p := range priorities {
		q, err := ch.QueueDeclare(
			"",    // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,
		)
		if err != nil {
			panic("Failed to declare queue")
		}
		err = ch.QueueBind(
			q.Name,           // queue name
			p,                // routing key
			"messages_topic", // exchange
			false, nil)
		if err != nil {
			panic("Failed to bind queue")
		}
		msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		qs[p] = msgs
	}

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	timeoutDuration := time.Duration(c.RabbitMessageTimeout) * time.Second
	var current *amqp.Delivery
	go func() {
		for {
			time.Sleep(timeoutDuration)

			current = nil
			select {
			default:
				break
			case msg := <-qs["high"]:
				current = &msg
			}
			if current == nil {
				select {
				case msg := <-qs["medium"]:
					current = &msg
				default:
					break
				}
			}
			if current == nil {
				select {
				case msg := <-qs["low"]:
					current = &msg
				default:
					break
				}
			}
			if current != nil {
				fmt.Println(current.Body)
				out := models.Message{}
				schema := models.SchemaMessage()
				err = avro.Unmarshal(schema, current.Body, &out)
				if err != nil {
					log.Fatal("Wrong schema")
				}
				cb(out)
			}
		}
	}()

	<-forever
}
