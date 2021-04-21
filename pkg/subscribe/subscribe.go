package subscribe

import (
	"errors"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	ErrTopicNotPresent = errors.New("Topic not found")
)

type Subscriber chan interface{}

type Publisher struct {
	topics map[string][]Subscriber
}

func NewPublisher() *Publisher {
	return &Publisher{
		make(map[string][]Subscriber),
	}
}

func (p *Publisher) HasTopic(topic string) bool {
	_, exist := p.topics[topic]
	return exist
}

func (p *Publisher) Publish(topic string, msg interface{}) error {
	if !p.HasTopic(topic) {
		return ErrTopicNotPresent
	}
	send := make(chan interface{})
	ok := make(chan struct{})
	go func() {
		msg := <-send
		ok <- struct{}{}
	}()
	count := 0
	for _, sub := range p.topics[topic] {
		timer := time.NewTimer(45 * time.Second)
		select {
		case <-timer.C:
			count++
		case <-ok:
		}
	}
	if count != 0 {
		log.Info().Msgf("Not every sub on topic \"%v\" accepted msg = %v", topic, msg)
	}
	return nil
}

func (p *Publisher) Subscribe(sub Subscriber, topic string) error {
	if p.HasTopic(topic) {
		p.topics[topic] = append(p.topics[topic], sub)
	} else {
		p.topics[topic] = []Subscriber{sub}
	}
	return nil
}
