package subscription

import (
	"errors"
	"quicksend/internal/user"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindActive(u *user.User) (*Subscription, error) {
	var sub Subscription

	err := r.db.
		Where("user_id = ? AND is_active = true AND end_at > ?", u.ID, time.Now().UTC()).
		First(&sub).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &sub, err
}

func (r *Repository) HasUsedTrial(u *user.User) (bool, error) {
	var count int64

	err := r.db.Model(&Subscription{}).
		Where("user_id = ? AND plan = ?", u.ID, PlanTrial).
		Count(&count).Error

	return count > 0, err
}

func (r *Repository) Create(sub *Subscription) error {
	return r.db.Create(sub).Error
}
