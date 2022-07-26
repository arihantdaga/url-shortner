package keygen

import (
	"context"
	"fmt"
	"sync/atomic"
	"url-shortner-keygen/lib"

	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
)

var MAX_KEYS int32 = 5
var KEYS_BUFFER = make([]string, MAX_KEYS)
var CURRENT_KEY_INDEX int32 = 0
var started = false

type KeygenService interface {
	GenerateKey(ctx context.Context, url string) string
	InsertKeys(ctx context.Context, n int)
}

type keyGenService struct {
	natsCon *nats.Conn
	pgCon   *pgx.Conn
}

func (k *keyGenService) GenerateKey(ctx context.Context, url string) string {
	// Technically i think it is not correct, we may need some mutex here,
	// since we are using the global array.
	// Also the database find and update operations have delay in between them, it may cause some ids to be recieved by multiple processes.
	var key string
	if CURRENT_KEY_INDEX >= MAX_KEYS || !started {
		started = true
		// fetch from DB again
		var shortkey string
		var id int
		var i int = 0
		var ids = make([]int32, MAX_KEYS)
		scanner, err := k.pgCon.Query(ctx, "SELECT id,shortkey FROM shortkeys WHERE taken = false LIMIT $1", MAX_KEYS)
		if err != nil {
			// fmt.Println(err)
			return ""
		}
		for scanner.Next() {
			scanner.Scan(&id, &shortkey)
			KEYS_BUFFER[i] = shortkey
			ids[i] = int32(id)
			i++
		}
		atomic.StoreInt32(&CURRENT_KEY_INDEX, 0)
		fmt.Println(KEYS_BUFFER)
		fmt.Println(ids)
		_, err = k.pgCon.Exec(ctx, "UPDATE shortkeys SET taken = true WHERE taken = false AND id = any($1)", ids)
		if err != nil {
			fmt.Println(err)
		}
		key = KEYS_BUFFER[CURRENT_KEY_INDEX]
		atomic.AddInt32(&CURRENT_KEY_INDEX, 1)
		// fmt.Println(key)
	} else {
		// fetch from buffer
		key = KEYS_BUFFER[CURRENT_KEY_INDEX]
		atomic.AddInt32(&CURRENT_KEY_INDEX, 1)
	}
	// key := lib.RandomStr(7)
	return key
}

func (k *keyGenService) InsertKeys(ctx context.Context, n int) {
	var i int
	var last_id int
	fmt.Println(k.pgCon.Config().Database)
	err := k.pgCon.QueryRow(ctx, "SELECT id FROM shortkeys ORDER BY id DESC LIMIT 1").Scan(&last_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(last_id)
	for i = last_id; i < last_id+n; i++ {
		fmt.Println(i + 1)
		key := lib.Base62(int64(i + 1000000000000))
		_, err := k.pgCon.Exec(ctx, "INSERT INTO shortkeys (id,shortkey,taken) VALUES ($1,$2,$3)", i+1, key, false)
		if err != nil {
			fmt.Println("Error inserting key: ", err)
		}
	}
}

func NewKeygenService(natsCon *nats.Conn, pgCon *pgx.Conn) KeygenService {
	return &keyGenService{
		natsCon: natsCon,
		pgCon:   pgCon,
	}
}
