package v1

import "github.com/google/uuid"

type Message struct {
	Key  string
	Body interface{}
}

func NewMessage(body any) *Message {
	return &Message{Key: uuid.NewString(), Body: body}
}
