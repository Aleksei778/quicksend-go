package subscription

import (
	"time"

	"quicksend/internal/user"

	"gorm.io/gorm"
)

type Plan string

const (
	PlanTrial    Plan = "trial"
	PlanStandard Plan = "standard"
	PlanPremium  Plan = "premium"
)

func (p Plan) DailyLimit() int {
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

func (p Plan) DaysCount() int {
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

	Plan                  Plan      `gorm:"not null"`
	IsActive              bool      `gorm:"default:true"`
	AutoRenew             bool      `gorm:"default:false"`
	StartedAt             time.Time `gorm:"not null"`
	EndAt                 time.Time `gorm:"not null"`
	CanceledAt            *time.Time
	FailedPaymentAttempts int       `gorm:"default:0"`
	UserID                uint      `gorm:"not null"`
	User                  user.User `gorm:"foreignKey:UserID"`
}
