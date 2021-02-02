package disperse

import (
	"fmt"
	"net"
)

type Connection struct {
	address   string
	consumers []Consumer
	Event     chan Event
}

func New(name string) Connection {
	newSocket := Connection{}
	newSocket.address = "/tmp/" + name + ".sock"
	newSocket.Event = make(chan Event)

	return newSocket
}

func (connection *Connection) setConsumer(consumerConnection net.Conn) {
	consumer := NewConsumer()
	consumer.connection = consumerConnection

	connection.consumers = append(connection.consumers, consumer)
	connection.Event <- Event{
		ConsumerId: consumer.id,
		Kind:       "connected",
	}
}

// getConsumerIndex get the position of a consumer off the consumer list
func (connection *Connection) getConsumerIndex(consumerId string) int {
	result := -1
	for index := range connection.consumers {
		if connection.consumers[index].id == consumerId {
			result = index
		}
	}

	return result
}

// killConsumer removes an active stopped/idle/active connection,
// and triggers a "disconnected" event
func (connection *Connection) killConsumer(consumerId string) {
	consumerPosition := connection.getConsumerIndex(consumerId)
	if consumerPosition > -1 {
		connection.consumers = append(connection.consumers[:consumerPosition], connection.consumers[consumerPosition+1:]...)
		connection.Event <- Event{
			ConsumerId: consumerId,
			Kind:       "disconnected",
		}
	}
}

// getConsumerMessage Get message of according consumer
func (connection *Connection) getConsumerMessage(consumerId string) {
	consumerPosition := connection.getConsumerIndex(consumerId)
	if consumerPosition > -1 {
		buf := make([]byte, 512)
		consumer := connection.consumers[consumerPosition]
		nr, err := consumer.connection.Read(buf)
		if err != nil {
			// return make([]byte, 0)
			fmt.Print("Errro reading message")
		}

		connection.Event <- Event{
			ConsumerId: consumerId,
			Kind:       "messafe",
			Data:       buf[0:nr],
		}
	}
}
