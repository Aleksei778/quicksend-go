package user

import (
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *Repository
}

func NewService(db *gorm.DB, r *Repository) *Service {
	return &Service{db: db, repository: r}
}

func (service *Service) FindOrCreate(dto FindOrCreate) (*User, error) {
	user, err := service.repository.FindByEmail(dto.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return service.create(dto)
	}

	return user, nil
}

func (service *Service) create(dto FindOrCreate) (*User, error) {
	user := &User{
		Email:      dto.Email,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		PictureUrl: dto.PictureUrl,
		OauthID:    dto.OauthID,
	}

	if err := service.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
