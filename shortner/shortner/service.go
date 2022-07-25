package shortner

import (
	"context"
	"errors"
	"time"
	"url-shortner/lib"

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
}

func (s *shortUrlService) GetUrl(ctx context.Context, shortkey string) (string, error) {
	// url, err := GetUrlFromDb(ctx, shortUrl)
	// if err != nil {
	// 	return "", err
	// }
	original_url := ""
	s.cassa.Query("SELECT original_url FROM shorturl WHERE shortkey = ?", shortkey).Scan(&original_url)
	return original_url, nil
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
func NewShortnerService(cassa *gocql.Session, natsCon *nats.Conn) ShortUrlService {
	s := shortUrlService{cassa: cassa, natsConn: natsCon}
	return &s
}
