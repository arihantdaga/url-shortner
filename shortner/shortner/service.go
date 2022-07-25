package shortner

import (
	"context"
	"errors"
	"fmt"
	"time"
	"url-shortner/lib"

	"github.com/go-redis/redis/v9"
	"github.com/gocql/gocql"
	"github.com/nats-io/nats.go"
)

type ShortUrlService interface {
	GetUrl(ctx context.Context, shortUrl string) (string, error)
	CreateUrl(ctx context.Context, shortUrl ShortUrl) (ShortUrl, error)
}

type shortUrlService struct {
	cassa    *gocql.Session
	natsConn *nats.Conn
	redis    *redis.Client
	cache    ShortnerCache
}

func (s *shortUrlService) GetUrl(ctx context.Context, shortkey string) (string, error) {
	original_url, err := s.cache.GetUrl(ctx, shortkey)
	if err == nil && original_url != "" {
		// fmt.Println("Cache hit")
		return original_url, nil
	} else {
		if err != nil {
			// fmt.Println(err.Error())
		}
		// fmt.Println("Cache miss")
		s.cassa.Query("SELECT original_url FROM shorturl WHERE shortkey = ?", shortkey).Scan(&original_url)
		if original_url != "" {
			err := s.cache.SetUrl(ctx, shortkey, original_url)
			if err != nil {
				fmt.Println(err.Error())
			}
			return original_url, nil
		} else {
			return "", errors.New("Not found")
		}
	}
}

func (s *shortUrlService) CreateUrl(ctx context.Context, shortUrl ShortUrl) (ShortUrl, error) {
	if shortUrl.Url == "" {
		return ShortUrl{}, errors.New("Url is empty")
	}
	shortUrl.CreatedAt = time.Now()
	shortKey := make(chan string)
	errChan := make(chan error)
	// go s.getShortKeyRandom(shortUrl.Url, shortKey)
	go s.getShortKeyFromKGS(shortUrl.Url, shortKey, errChan)
	select {
	case msg1 := <-shortKey:
		shortUrl.Shortkey = msg1
	case err := <-errChan:
		return ShortUrl{}, err
	}
	// shortUrl.Shortkey = fmt.Sprintf("%s", <-shortKey)
	if err := s.cassa.Query("INSERT INTO shorturl (shortkey, original_url, CreatedAt) VALUES (?, ?, ?)", shortUrl.Shortkey, shortUrl.Url, shortUrl.CreatedAt).Exec(); err != nil {
		return ShortUrl{}, err
	} else {
		return shortUrl, nil
	}
}

func (s *shortUrlService) getShortKeyRandom(url string, shortKey chan string) {
	// time.Sleep(2 * time.Second)
	shortKey <- lib.RandomStr(7)
}
func (s *shortUrlService) getShortKeyFromKGS(url string, shortKey chan string, errChan chan error) {
	data := url
	msg, err := s.natsConn.Request("shortkey", []byte(data), time.Second*3)
	if err != nil {
		errChan <- err
	} else {
		shortKey <- string(msg.Data)
	}
}
func NewShortnerService(cassa *gocql.Session, natsCon *nats.Conn, redis *redis.Client) ShortUrlService {
	cache := NewShortnerCache(redis)
	s := shortUrlService{cassa: cassa, natsConn: natsCon, redis: redis, cache: cache}
	return &s
}
