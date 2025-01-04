package v1

import (
	"fmt"

	"github.com/vnkdj5/low-level-design/pub-sub-system/utils"
)

type Publisher struct {
	Topics *utils.Set[*Topic]
}

func NewPublisher() *Publisher {
	return &Publisher{
		Topics: utils.NewSet[*Topic](),
	}
}

func (p *Publisher) RegisterTopic(topic *Topic) {
	p.Topics.Add(topic)
}

func (p *Publisher) PublishMessage(topic *Topic, msg *Message) {
	exists := p.Topics.Contains(topic)
	if !exists {
		fmt.Println(fmt.Sprintf("Publisher cannot publish to the topic: %s", topic.Name))
	}
	topic.PublishMessageSync(msg)
}
