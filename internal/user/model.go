package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string `gorm:"uniqueIndex;not null"`
	PictureUrl string
	OauthID    string `gorm:"not null"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
}
