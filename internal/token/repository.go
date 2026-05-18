package token

import (
	"errors"
	"quicksend/internal/user"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUser(user user.User) (*Token, error) {
	var token Token
	err := r.db.Where("user_id = ?", user.ID).First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &token, err
}
