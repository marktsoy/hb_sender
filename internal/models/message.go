package models

import (
	"log"

	"github.com/hamba/avro"
)

// Message ...
type Message struct {
	Content  string `json:"string" avro:"content"`
	Priority int    `json:"priority" avro:"priority"`
	BundleID string `json:"bundle_id" avro:"bundle_id"`
}

var (
	// MessageSchema ...
	MessageSchema *avro.RecordSchema
)

// SchemaMessage ...
func SchemaMessage() *avro.RecordSchema {
	if MessageSchema == nil {
		schema, err := avro.Parse(`{
			"type": "record",
			"name": "message",
			"fields" : [
				{"name": "content", "type": "string"},
				{"name": "status", "type": "int"},
				{"name": "bundle_id", "type": "string"}
			]
		}`)
		if err != nil {
			log.Fatal(err)
		}
		MessageSchema = schema.(*avro.RecordSchema)
	}
	return MessageSchema
}
