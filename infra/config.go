package infra

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Cfg is a configuration variable for the app.
var Cfg Config
var once sync.Once

// Config is a general configuration.
type Config struct {
	Mode     string
	Port     string
	Debug    string
	Postgres PostgresConfig
	Redis    RedisConfig
	Casbin   CasbinConfig
}

// LoadEnv loads configuration from env variables.
func LoadEnv() {
	if os.Getenv("APP_ENV") == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	once.Do(func() {
		Cfg = Config{
			Mode:  os.Getenv("MODE"),
			Port:  os.Getenv("PORT"),
			Debug: os.Getenv("DEBUG"),
			Postgres: PostgresConfig{
				Host:     os.Getenv("POSTGRES_HOST"),
				Port:     os.Getenv("POSTGRES_PORT"),
				User:     os.Getenv("POSTGRES_USER"),
				Password: os.Getenv("POSTGRES_PASS"),
				DBName:   os.Getenv("POSTGRES_DBNAME"),
				SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
				SSLCert:  os.Getenv("POSTGRES_SSL_CERT"),
				SSLKey:   os.Getenv("POSTGRES_SSL_KEY"),
			},
			Redis: RedisConfig{
				Host:     os.Getenv("REDIS_HOST"),
				Port:     os.Getenv("REDIS_PORT"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       redisDB(),
			},
			Casbin: CasbinConfig{
				ModelPath:  os.Getenv("CASBIN_MODEL_PATH"),
				PolicyPath: os.Getenv("CASBIN_POLICY_PATH"),
			},
		}
	})
}

// func newNatsConn() (*nats.Conn, *nats.EncodedConn, error) {
// 	nc, err := nats.Connect("localhost:4222")
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return nc, ec, nil
// }
