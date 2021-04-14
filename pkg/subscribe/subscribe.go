package subscribe

import "errors"

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
	for _, sub := range p.topics[topic] {
		sub <- msg
	}
	return nil
}

func (p *Publisher) Subsciribe(sub Subscriber, topic string) error {
	if p.HasTopic(topic) {
		p.topics[topic] = append(p.topics[topic], sub)
	} else {
		p.topics[topic] = []Subscriber{sub}
	}
	return nil
}

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	var once sync.Once

// 	for i := 0; i < 1000; i++ {
// 		once.Do(func() {
// 			fmt.Println(" *Executed *")
// 		})
// 	}
// }
