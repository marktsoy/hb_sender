package models

import (
	"log"

	"github.com/hamba/avro"
)

// Subscription ...
type Subscription struct {
	ChatID       int64 `json:"chat_id" avro:"chat_id"`
	IsSubscribed bool  `json:"is_subscribed" avro:"is_subscribed"`
}

var (
	// SubscriptionSchema ...
	SubscriptionSchema *avro.RecordSchema
)

// SchemaSubscription ...
func SchemaSubscription() *avro.RecordSchema {
	if SubscriptionSchema == nil {
		schema, err := avro.Parse(`{
			"type": "record",
			"name": "subscription",
			"fields" : [
				{"name": "chat_id", "type": "long"},
				{"name": "is_subscribed", "type": "boolean"}
			]
		}`)
		if err != nil {
			log.Fatal(err)
		}
		SubscriptionSchema = schema.(*avro.RecordSchema)
	}
	return SubscriptionSchema
}
