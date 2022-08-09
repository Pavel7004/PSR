package subscribe

import (
	"reflect"
	"testing"
)

func TestNewPublisher(t *testing.T) {
	tests := []struct {
		name string
		want *Publisher
	}{
		{
			name: "Test creating new publisher",
			want: &Publisher{make(map[string][]ISubscriber)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPublisher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPublisher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublisher_Publish(t *testing.T) {
	type fields struct {
		topics map[string][]ISubscriber
	}
	type args struct {
		topic string
		msg   interface{}
	}
	testSubscribers := map[string][]ISubscriber{
		"Topic not exist": {},
		"Topic exists, msg received": {
			NewSubscriber(0),
		},
		"Topic exists, two subs, msg received": {
			NewSubscriber(0),
			NewSubscriber(0),
		},
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected string
		wantErr  bool
	}{
		{
			name:     "Topic exists, msg received",
			fields:   fields{map[string][]ISubscriber{"topic 1": testSubscribers["Topic exists, msg received"]}},
			args:     args{"topic 1", "msg"},
			expected: "msg",
			wantErr:  false,
		},
		{
			name: "Topic exists, two subs, msg received",
			fields: fields{
				map[string][]ISubscriber{
					"topic 1": testSubscribers["Topic exists, two subs, msg received"],
				},
			},
			args:     args{"topic 1", "msg"},
			expected: "msg",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Publisher{
				topics: tt.fields.topics,
			}
			ready := make(chan struct{})
			go func() {
				for index, sub := range testSubscribers[tt.name] {
					if msg := sub.Receive(); tt.expected != msg {
						t.Errorf(`Publisher.Publish() Gotten wrong message in subscriber #%d
								expected message=%v
								actual %v`,
							index, tt.expected, msg)
					}
				}
				ready <- struct{}{}
			}()
			if err := p.Publish(tt.args.topic, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Publisher.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
			<-ready
		})
	}
}

func TestPublisher_Subscribe(t *testing.T) {
	type fields struct {
		topics map[string][]ISubscriber
	}
	type args struct {
		sub   *Subscriber
		topic string
	}
	addingSubscriber := NewSubscriber(0)
	existingSubscriber := NewSubscriber(0)
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected []ISubscriber
		wantErr  bool
	}{
		{
			name:     "topic exist",
			fields:   fields{topics: map[string][]ISubscriber{"topic 1": {}, "topic 2": {existingSubscriber}}},
			args:     args{addingSubscriber, "topic 1"},
			expected: []ISubscriber{addingSubscriber},
			wantErr:  false,
		},
		{
			name:     "topic exist, second subscriber",
			fields:   fields{topics: map[string][]ISubscriber{"topic 1": {}, "topic 2": {existingSubscriber}}},
			args:     args{addingSubscriber, "topic 2"},
			expected: []ISubscriber{existingSubscriber, addingSubscriber},
			wantErr:  false,
		},
		{
			name:     "topic is not exist, one subscriber",
			fields:   fields{topics: map[string][]ISubscriber{"topic 1": {}, "topic 2": {existingSubscriber}}},
			args:     args{addingSubscriber, "topic 3"},
			expected: []ISubscriber{addingSubscriber},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Publisher{
				topics: tt.fields.topics,
			}
			if err := p.Subscribe(tt.args.sub, tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("Publisher.Subsciribe() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.expected, p.topics[tt.args.topic]) {
				t.Errorf("Publisher.Subsciribe() lists of subscribers not equal. Expected: %v, actual: %v",
					tt.expected, p.topics[tt.args.topic])
			}
		})
	}
}
