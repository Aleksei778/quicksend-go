package models

import (
	"time"

	"quicksend/internal/user"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	gorm.Model
	UserID            uint          `gorm:"not null"`
	SubscriptionID    uint          `gorm:"not null"`
	ExternalPaymentID string        `gorm:"uniqueIndex"`
	Amount            float64       `gorm:"type:decimal(10,2)"`
	Currency          string        `gorm:"default:'RUB'"`
	Status            PaymentStatus `gorm:"default:'pending'"`
	PaymentMethod     string
	Description       string
	PaidAt            *time.Time
	User              user.User    `gorm:"foreignKey:UserID"`
	Subscription      Subscription `gorm:"foreignKey:SubscriptionID"`
}
