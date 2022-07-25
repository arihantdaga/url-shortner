package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

func ReadConfig(envFile string) error {
	// Load JSON config.
	if err := k.Load(file.Provider(envFile), dotenv.Parser()); err != nil {
		log.Println("Error reading config file :", err)
	} else {
		log.Println("Config file loaded successfully")
	}
	if err := k.Load(env.Provider("TODAPP_", ".", func(s string) string {
		news := strings.TrimPrefix(s, "TODAPP_")
		return news
	}), nil); err != nil {
		return err
	}
	return nil
}

func Get(key string) interface{} {
	return k.Get(key)
}
