package token

import (
	"quicksend/internal/config"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	config     *config.Config
	repository *Repository
}

func NewService(db *gorm.DB, r *Repository, cfg *config.Config) *Service {
	return &Service{db: db, repository: r, config: cfg}
}

func (service *Service) FindOrCreate(dto FindOrCreate) (*Token, error) {
	token, err := service.repository.FindByUser(dto.user)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return service.create(dto)
	}

	return token, nil
}

func (service *Service) create(dto FindOrCreate) (*Token, error) {
	token := &Token{
		UserID:  dto.user.ID,
		Access:  dto.access,
		Refresh: dto.refresh,
		Expiry:  dto.expiry,
	}

	if err := service.db.Create(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (service *Service) update(access string, expiry time.Time, token *Token) (*Token, error) {
	token.Access = access
	token.Expiry = expiry

	if err := service.db.Save(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}
