package subscribe

import (
	"errors"

	"github.com/rs/zerolog/log"
)

var (
	ErrTopicNotPresent = errors.New("Topic not found")
)

type Publisher struct {
	topics map[string][]ISubscriber
}

func NewPublisher() *Publisher {
	return &Publisher{
		make(map[string][]ISubscriber),
	}
}

func (p *Publisher) Publish(topic string, msg interface{}) error {
	for _, sub := range p.topics[topic] {
		if err := sub.Send(msg); err != nil {
			log.Error().Err(err).Msgf("publish error, topic: %s, msg: %v", topic, msg)
			return err
		}
	}

	return nil
}

func (p *Publisher) Subscribe(sub ISubscriber, topic string) error {
	_, exist := p.topics[topic]
	if exist {
		p.topics[topic] = append(p.topics[topic], sub)
	} else {
		p.topics[topic] = []ISubscriber{sub}
	}
	return nil
}
