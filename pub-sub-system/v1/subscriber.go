package v1

import (
	"fmt"
)

type Subscriber interface {
	OnMessage(msg *Message) error
}

type ConsoleSubscriber struct {
	Name string
}

func NewConsoleSubscriber(name string) Subscriber {
	return &ConsoleSubscriber{Name: name}
}

func (cs *ConsoleSubscriber) OnMessage(msg *Message) error {

	if msg == nil {
		return fmt.Errorf("message is nil, not consuming")
	}
	fmt.Printf("Subscriber: %s \tMessageKey: %s\tMessageBody: %v\n", cs.Name, msg.Key, msg.Body)
	return nil
}
