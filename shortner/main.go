package main

import (
	"fmt"
	"log"
	"os"
	config "url-shortner/config"
	"url-shortner/database"
	"url-shortner/queue"
	shortner "url-shortner/shortner"

	"github.com/gofiber/fiber/v2"
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
	port := config.Get("APP_PORT").(string)
	cassaUri := config.Get("CASSA_URI").(string)

	cassSess, err := database.NewCasDb(cassaUri)

	// dbUri := config.Get("DB_URI").(string)
	// dbClient, err := database.NewDb(dbUri)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	} else {
		defer cassSess.Close()
	}

	natsUri := config.Get("NATS_URI").(string)
	natsConn, err := queue.NewNats(natsUri)
	if err != nil {
		log.Println("Error connecting to queue:", err)
		return err
	} else {
		defer natsConn.Close()
	}

	app := fiber.New()
	api := app.Group("/api")
	// authRoute := app.Group("/auth")
	apiv1 := api.Group("/v1")
	shortner.Routes(apiv1, shortner.NewShortnerService(cassSess, natsConn))

	// Setup Auth
	// auth.Routes(authRoute, auth.NewAuthService(dbClient))
	iface := fmt.Sprintf(":%s", port)
	log.Println(iface)
	return app.Listen(iface)
}
