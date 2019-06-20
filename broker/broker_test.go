package broker

import (
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

func TestLocalIP(t *testing.T) {
	LocalAddress()
}

func natsErrHandler(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
	fmt.Printf("error: %v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			fmt.Printf("couldn't get pending messages: %v", err)
			return
		}
		fmt.Printf("Falling behind with %d pending messages on subject %q.\n",
			pendingMsgs, sub.Subject)
		// Log error, notify operations...
	}
	// check for other errors
}

func TestNatsSubBreak(t *testing.T) {
	ns1, err := nats.Connect("nats://127.0.0.1:4222", nats.ErrorHandler(natsErrHandler))
	assert.Nil(t, err)
	s, err := ns1.Subscribe("register", func(msg *nats.Msg) {
		fmt.Printf("register:%v\n", string(msg.Data))
	})
	defer s.Unsubscribe()
	assert.Nil(t, err)

	ns1.Status()
	ns2, _ := nats.Connect("nats://127.0.0.1:4222")
	ns2.Subscribe("client2", func(msg *nats.Msg) {
		fmt.Printf("ping client2:%v\n", string(msg.Data))
	})
	ns2.Publish("register", []byte("hello"))

	ns1.Publish("client2", []byte("ping"))
	time.Sleep(2 * time.Second)
	ns2.Close()
	ns1.Publish("client2", []byte("ping"))
	fmt.Println("close ns2")
	time.Sleep(2 * time.Second)
	ns1.Publish("client2", []byte("ping"))
	time.Sleep(2 * time.Second)
}
