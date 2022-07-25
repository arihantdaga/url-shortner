package database

import (
	"log"
	"time"

	gocql "github.com/gocql/gocql"
)

func NewCasDb(uri string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(uri)
	cluster.Keyspace = "shorturl"
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Timeout = 10 * time.Second
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Cassandra")
	return session, nil
	// return nil, nil
}
