package subscribe

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
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
	send := make(chan struct {
		sub *Subscriber
		msg interface{}
	})
	ok := make(chan struct{})
	go func() {
		for {
			directions := <-send
			if directions.sub == nil {
				break
			}
			(*directions.sub) <- directions.msg
			ok <- struct{}{}
		}
	}()
	count := 0
	for _, sub := range p.topics[topic] {
		timer := time.NewTimer(45 * time.Second)
		send <- struct {
			sub *Subscriber
			msg interface{}
		}{
			&sub,
			msg,
		}
		select {
		case <-timer.C:
			count++
		case <-ok:
		}
	}
	send <- struct {
		sub *Subscriber
		msg interface{}
	}{
		nil,
		nil,
	}
	if count != 0 {
		log.Info().Msgf("%d subs on topic \"%v\" accepted msg = %v", count, topic, msg)
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
