package app_test

import (
	"testing"
	"time"

	"github.com/streadway/amqp"

	"github.com/hamba/avro"
	"github.com/marktsoy/hb_sender/internal/app"
	"github.com/marktsoy/hb_sender/internal/models"
	"github.com/stretchr/testify/assert"
)

func testSender(t *testing.T) (*amqp.Connection, *amqp.Channel) {
	t.Helper()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to open a channel")
	}

	if err != nil {
		panic(err)
	}
	return conn, ch
}

func TestMessageListener_test(t *testing.T) {
	conn, ch := testSender(t)
	consumedMsg := []models.Message{}
	msgs := []models.Message{
		{ID: 1, Content: "First message", Status: 0, Priority: 0, BundleID: 1},
		{ID: 2, Content: "Second message", Status: 0, Priority: 4, BundleID: 3},
		{ID: 3, Content: "Third message", Status: 0, Priority: 0, BundleID: 2},
	}
	defer func() {
		conn.Close()
		ch.Close()
	}()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}

	go func() {
		t.Log("Starting Listener")
		app.ListenToMessages("hello", func(m models.Message) {
			t.Log(map[string]interface{}{
				"id":       m.ID,
				"content":  m.Content,
				"priority": m.Priority,
				"status":   m.Status,
				"bundle":   m.BundleID,
			})
			consumedMsg = append(consumedMsg, m)
		})
		t.Log("Listener closed")
	}()

	for _, msg := range msgs {

		t.Logf("Sending to queue msg: %v", msg.ID)
		schema := models.SchemaMessage()

		data, err := avro.Marshal(schema, msg)
		assert.NoError(t, err)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        data,
			})
		assert.NoError(t, err)
	}

	time.Sleep(5 * time.Second)
	assert.Equal(t, len(msgs), len(consumedMsg))
}

func TestMessage_PopulateQueue(t *testing.T) {
	conn, ch := testSender(t)
	msgs := []models.Message{
		{ID: 1, Content: "First message", Status: 0, Priority: 0, BundleID: 1},
		{ID: 2, Content: "Second message", Status: 0, Priority: 4, BundleID: 3},
		{ID: 3, Content: "Third message", Status: 0, Priority: 0, BundleID: 2},
	}
	defer func() {
		conn.Close()
		ch.Close()
	}()

	q, err := ch.QueueDeclare(
		"messages", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		panic(err)
	}

	for _, msg := range msgs {

		t.Logf("Sending to queue msg: %v", msg.ID)
		schema := models.SchemaMessage()

		data, err := avro.Marshal(schema, msg)
		assert.NoError(t, err)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        data,
			})
		assert.NoError(t, err)
	}

	time.Sleep(5 * time.Second)
}
