package db

import (
	"fmt"
	"quicksend/internal/subscription"

	"quicksend/internal/config"
	"quicksend/internal/models"
	"quicksend/internal/token"
	"quicksend/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&user.User{},
		&token.Token{},
		&models.Campaign{},
		&models.Recipient{},
		&models.Attachment{},
		&subscription.Subscription{},
		&models.Payment{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}
