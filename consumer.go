package disperse

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"net"
	"strconv"
	"time"
)

type Consumer struct {
	id         string
	connection net.Conn
}

// computeId generates an unique identifier
// in order to tag each connected consumer
func (consumer *Consumer) setConsumerId() {
	w := sha256.New()
	t := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	w.Write([]byte(t + ""))

	buf := bytes.NewBuffer(nil)
	bw := base64.NewEncoder(base64.StdEncoding, buf)
	bw.Write(w.Sum(nil))
	bw.Close()

	consumer.id = buf.String()
}

// New setups a consumer
func NewConsumer() Consumer {
	// Computes an unique Id
	consumer := Consumer{}
	consumer.setConsumerId()

	return consumer
}
