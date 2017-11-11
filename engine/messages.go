package engine

import "time"

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
// of the action, any error that might of occured and how long the action took to process.
type Reply struct {
	Duration time.Duration
	Data     interface{}
	Error    error
}
