package ds

import "time"

type MassRequest struct {
	ID                uint64              `gorm:"primaryKey"`
	Status            uint8               `gorm:"not null;default:1"`  // 1 - draft, 2 - deleted, 3 - pending, 4 - resolved, 5 - rejected
	UserID            uint64              `gorm:"not null"`
	User              User                `gorm:"foreignKey:UserID;references:ID"`
	ModeratorId       uint64              `gorm:"default:null"`
	Moderator         User                `gorm:"foreignKey:ModeratorId;references:ID"`
	MassRequestToClass []MassRequestToClass `gorm:"foreignKey:RequestID"`
	CreatedAt         time.Time           `gorm:"not null;default:now()"`
	FormedAt          time.Time           `gorm:"default:null"`
	ClosedAt          time.Time           `gorm:"default:null"`
}