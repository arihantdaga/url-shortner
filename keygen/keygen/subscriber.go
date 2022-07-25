package keygen

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

func Subscriber(nc *nats.Conn, service KeygenService) {
	// Replies
	nc.Subscribe("shortkey", func(m *nats.Msg) {
		fmt.Println("Received Request ", string(m.Data))
		shortkey := service.GenerateKey(context.Background(), string(m.Data))
		nc.Publish(m.Reply, []byte(shortkey))
	})
}
