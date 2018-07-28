package engine

import (
	"time"
)

// NewMessage creates a new Message to be sent to the Engine for evaluation.
func NewMessage(action interface{}) *Message {
	return &Message{
		Action: action,
		Reply:  make(chan *Reply),
	}
}

// Message is a payload that contains an operation that the Engine can process and channel that must be received on by
// the goroutine sending the Message.
type Message struct {
	Action interface{}
	Reply  chan *Reply
}

// Reply is a payload that the Engine sends in response to a Message. It contains any resulting data from its processing
// of the action, any error that might of occurred and how long the action took to process.
type Reply struct {
	Duration time.Duration
	Data     interface{}
	Error    error
}

// MessageChannel is abstraction of a channel that handles engine messages. This provides us a means of implementing
// slightly more strict synchronization behavior during testing.
type MessageChannel interface {
	Receive() *Message
	Send(*Message) error
	Close()
}

func newMessageChannel() messageChannel {
	return messageChannel{make(chan *Message, 100)}
}

type messageChannel struct {
	messages chan *Message
}

func (b messageChannel) Receive() *Message {
	select {
	case msg := <-b.messages:
		return msg
	default:
		return nil
	}
}

func (b messageChannel) Send(msg *Message) error {
	b.messages <- msg
	return nil
}
func (b messageChannel) Close() { close(b.messages) }
