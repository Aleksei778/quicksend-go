package main

import (
	"log"
	"quicksend/internal/config"
	"quicksend/internal/crypto"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: cfg.BuggregatorDSN,
	})
	if err != nil {
		log.Fatalf("failed to init sentry: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	crypto.Init(cfg.EncryptionKey)

}
