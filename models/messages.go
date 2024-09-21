package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
)

type Message struct {
	ID             gocql.UUID `json:"id"`
	ConversationID gocql.UUID `json:"conversation_id"`
	SenderId       gocql.UUID `json:"sender_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Body           string     `json:"body"`
	IsSoftDeleted  bool       `json:"is_soft_deleted"`
}

var messageMetadata = table.Metadata{
	Name: "messaging_keyspace.messages",
	Columns: []string{
		"id",              //id for the message 👌🏼
		"conversation_id", //id for the conversation
		"sender_id",       //id for the sender
		"created_at",      //time when the message was created 👌🏼
		"updated_at",      //time when the message
		"body",            //body of the message
		"is_soft_deleted", //whether the message is soft deleted or not
	},
	PartKey: []string{"id"},
	SortKey: []string{"conversation_id"},
}

var MessageTable = table.New(messageMetadata)
