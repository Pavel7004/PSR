package subscribe

import "errors"

// Broker

// type Subscriber struct {
// 	name string
// }

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

// broker.Publish("topic1", channel)

// broker.Subscribe(vasya, "topic1")

func (p *Publisher) Publish(topic string, msg interface{}) error {
	if _, has := p.topics[topic]; !has {
		return ErrTopicNotPresent
	}
	for _, sub := range p.topics[topic] {
		sub <- msg
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
