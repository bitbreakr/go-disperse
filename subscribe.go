package disperse

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ping sends a ping event to consumer,
// disconnect a consumer if ping fails
func (connection *Connection) ping() {
	for {
		time.Sleep(1 * time.Second)
		if len(connection.consumers) > 0 {
			//do some ping, if no response then kill it
			for _, consumer := range connection.consumers {
				_, pingError := consumer.connection.Write([]byte("hunga"))
				if pingError != nil {
					// fmt.Print("PING ERROR")
					connection.killConsumer(consumer.id)
				} else {
					connection.getConsumerMessage(consumer.id)
				}
			}
		}
	}
}

func (connection *Connection) Subscribe() {
	listener, listenError := net.Listen("unix", connection.address)
	if listenError != nil {
		log.Fatal("Listen error: ", listenError)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(listener net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		listener.Close()
		os.Exit(0)
	}(listener, sigc)

	go connection.ping()

	for {
		consumerConnection, consumerConnectionError := listener.Accept()
		if consumerConnectionError != nil {
			log.Fatal("Accept error: ", consumerConnectionError)
		}

		// Sets consumer
		connection.setConsumer(consumerConnection)
	}
}
