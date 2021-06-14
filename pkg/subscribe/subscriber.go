package subscribe

type Subscriber struct {
	msgCh chan interface{}
}

func NewSubscriber(bufferSize int) *Subscriber {
	return &Subscriber{
		make(chan interface{}, bufferSize),
	}
}

func (s *Subscriber) Send(msg interface{}) error {
	s.msgCh <- msg
	return nil
}

func (s *Subscriber) Receive() interface{} {
	return <-s.msgCh
}
