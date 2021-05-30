package messagebroker

import (
	"github.com/nats-io/nats.go"
)

// MessageBroker is the interface to publish, consume messages
type MessageBroker interface {
	// Publish publishes message to the topic
	Publish(topic string, message ...Message) error
	// Subscribe subscribes to a chanel mesagges of topic
	Subscribe(topic string) (chan Message, error)
	// Close end MessageBroker
	Close()
}

type ncBroker struct {
	conn *nats.Conn
}

// Message message being sent or received
type Message struct {
	Payload []byte
}

// Config configuration to start MessageBroker
type Config struct {
	URL string
}

// New returns a new MessageBroker
func New(c Config) (MessageBroker, error) {
	nc, err := nats.Connect(c.URL)
	if err != nil {
		return nil, err
	}
	return &ncBroker{conn: nc}, nil
}

func (mb ncBroker) Publish(topic string, message ...Message) error {
	if err := mb.conn.Publish(topic, message[0].Payload); err != nil {
		return err
	}
	return nil
}

func (mb ncBroker) Subscribe(topic string) (chan Message, error) {
	msgChan := make(chan Message)
	mb.conn.Subscribe(topic, func(m *nats.Msg) {
		msgChan <- Message{Payload: m.Data}
	})
	return msgChan, nil
}

func (mb ncBroker) Close() {
	mb.conn.Close()
}
