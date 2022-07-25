package queue

import (
	"github.com/nats-io/nats.go"
)

func NewNats(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url)
	return nc, err
}
