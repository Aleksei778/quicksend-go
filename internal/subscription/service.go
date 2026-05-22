package subscription

import
import (
	"fmt"
	"quicksend/internal/user"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrial(u *user.User) error {
	used, err := s.repo.HasUsedTrial(u)
	if err != nil {
		return fmt.Errorf("subscription: check trial: %w", err)
	}
	if used {
		return nil
	}

	now := time.Now().UTC()
	sub := &Subscription{
		UserID: u.ID,
		Plan: PlanTrial,
		IsActive: true,
		StartedAt: now,
		EndAt: now.AddDate(0, 0, ),
	}
}
