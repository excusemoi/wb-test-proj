package config

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	"log"
	"os"
	"sync"
)

// Config is a config :)
type Config struct {
	PgUrl					string `env:"PG_URL"`
	PgPort                	string `env:"PG_PORT"`
	PgHostName				string `env:"PG_HOSTNAME"`
	PgDb                 	string `env:"PG_DB"`
	PgUser              	string `env:"PG_USER"`
	PgPassword           	string `env:"PG_PASSWORD"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment
func Get(path string) *Config {
	once.Do(func() {
		envType := os.Getenv("ENV")
		if envType == "" {
			envType = "dev"
		}
		if err := configor.New(&configor.Config{Environment: envType}).Load(&config, path); err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", " \t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})
	return &config
}