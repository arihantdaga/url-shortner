package keygen

import (
	"context"
	"url-shortner-keygen/lib"

	"github.com/nats-io/nats.go"
)

type KeygenService interface {
	GenerateKey(ctx context.Context, url string) string
}

type keyGenService struct {
	natsCon *nats.Conn
}

func (k *keyGenService) GenerateKey(ctx context.Context, url string) string {
	key := lib.RandomStr(7)
	return key
}

func NewKeygenService(natsCon *nats.Conn) KeygenService {
	return &keyGenService{
		natsCon: natsCon,
	}
}
