package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	config "url-shortner-keygen/config"
	"url-shortner-keygen/database"
	keygen "url-shortner-keygen/keygen"
	queue "url-shortner-keygen/queue"
)

var CONFIG_FILE = ".env"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Running application in %v environment", os.Getenv("APP_ENV"))

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	if err := config.ReadConfig(CONFIG_FILE); err != nil {
		return err
	}
	natsUri := config.Get("NATS_URI").(string)
	natsConn, err := queue.NewNats(natsUri)
	if err != nil {
		log.Println("Error connecting to queue:", err)
		return err
	} else {
		defer natsConn.Close()
	}
	// cassaUri := config.Get("CASSA_URI").(string)
	// cassSess, err := database.NewCasDb(cassaUri)
	// if err != nil {
	// 	log.Println("Error connecting to database:", err)
	// 	return err
	// } else {
	// 	defer cassSess.Close()
	// }
	pgUri := config.Get("PG_URI").(string)
	pgCon, err := database.NewPgDb(pgUri)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	} else {
		defer pgCon.Close(context.Background())
	}

	keygenService := keygen.NewKeygenService(natsConn, pgCon)
	// keygenService.InsertKeys(context.TODO(), 20000)
	keygen.Subscriber(natsConn, keygenService)

	// Run forever
	go forever()
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	return nil
}
func forever() {
	for {
		// fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
