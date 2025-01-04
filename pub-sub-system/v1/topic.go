package v1

import (
	"sync"

	"github.com/vnkdj5/low-level-design/pub-sub-system/utils"
)

type Topic struct {
	Name        string
	Subscribers *utils.Set[Subscriber]
	mu          sync.RWMutex
}

func NewTopic(name string) *Topic {
	return &Topic{Name: name,
		Subscribers: utils.NewSet[Subscriber](),
		mu:          sync.RWMutex{},
	}
}

func (t *Topic) AddSubscriber(s Subscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Subscribers.Add(s)
}

func (t *Topic) RemoveSubscriber(s Subscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Subscribers.Remove(s)
}

func (t *Topic) PublishMessageSync(msg *Message) {
	if msg == nil {
		return
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, sub := range t.Subscribers.Items() {
		sub.OnMessage(msg)
	}

}
