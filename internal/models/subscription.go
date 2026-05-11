package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionPlan string

const (
	PlanTrial    SubscriptionPlan = "trial"
	PlanStandard SubscriptionPlan = "standard"
	PlanPremium  SubscriptionPlan = "premium"
)

func (p SubscriptionPlan) DailyLimit() int {
	switch p {
	case PlanTrial:
		return 50
	case PlanStandard:
		return 500
	case PlanPremium:
		return 2000
	default:
		return 0
	}
}

func (p SubscriptionPlan) DaysCount() int {
	switch p {
	case PlanTrial:
		return 10
	default:
		now := time.Now()
		return time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	}
}

type Subscription struct {
	gorm.Model

	Plan                  SubscriptionPlan `gorm:"not null"`
	IsActive              bool             `gorm:"default:true"`
	AutoRenew             bool             `gorm:"default:false"`
	StartedAt             time.Time        `gorm:"not null"`
	EndAt                 time.Time        `gorm:"not null"`
	CanceledAt            *time.Time
	FailedPaymentAttempts int  `gorm:"default:0"`
	UserID                uint `gorm:"not null"`
	User                  User `gorm:"foreignKey:UserID"`
}
