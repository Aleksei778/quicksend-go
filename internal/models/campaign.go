package models

import (
	"time"

	"gorm.io/gorm"
)

type CampaignStatus string

const (
	StatusDraft     CampaignStatus = "draft"
	StatusScheduled CampaignStatus = "scheduled"
	StatusSending   CampaignStatus = "sending"
	StatusSent      CampaignStatus = "sent"
	StatusPaused    CampaignStatus = "paused"
	StatusCancelled CampaignStatus = "cancelled"
)

type Campaign struct {
	gorm.Model
	SenderName   string
	Subject      string         `gorm:"not null"`
	BodyTemplate string         `gorm:"type:text;not null"`
	Status       CampaignStatus `gorm:"default:'draft'"`
	StartedAt    *time.Time
	EndAt        *time.Time
	Timezone     string       `gorm:"default:'UTC'"`
	UserID       uint         `gorm:"not null"`
	User         User         `gorm:"foreignKey:UserID"`
	Recipients   []Recipient  `gorm:"foreignKey:CampaignID"`
	Attachments  []Attachment `gorm:"foreignKey:CampaignID"`
}

type Recipient struct {
	gorm.Model
	Email      string `gorm:"not null;index"`
	SentAt     *time.Time
	CampaignID uint `gorm:"not null"`
}

type Attachment struct {
	gorm.Model
	Filename   string `gorm:"not null"`
	Size       int64
	Mimetype   string `gorm:"default:'application/octet-stream'"`
	Content    string `gorm:"type:text;not null"` // base64
	CampaignID uint   `gorm:"not null"`
}
