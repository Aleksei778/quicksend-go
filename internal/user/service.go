package user

import (
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *Repository
}

func NewService(db *gorm.DB, r *Repository) *Service {
	return &Service{db: db, repo: r}
}

func (svc *Service) FindOrCreate(dto FindOrCreate) (*User, error) {
	user, err := svc.repo.FindByEmail(dto.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return svc.create(dto)
	}

	return user, nil
}

func (svc *Service) FindByID(id uint) (*User, error) {
	return svc.repo.FindByID(id)
}

func (svc *Service) FindByEmail(email string) (*User, error) {
	return svc.repo.FindByEmail(email)
}

func (svc *Service) create(dto FindOrCreate) (*User, error) {
	user := &User{
		Email:      dto.Email,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		PictureUrl: dto.PictureUrl,
		OauthID:    dto.OauthID,
	}

	if err := svc.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
