package subscribe

type ISubscriber interface {
	Send(interface{}) error
	Receive() interface{}
}
