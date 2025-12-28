package ds

import (
	"time"
)

type MassRequest struct {
	ID                 uint64               `gorm:"primaryKey" json:"id"`
	Status             uint8                `json:"status"`
	UserID             uint64               `json:"user_id"`
	ModeratorId        *uint64              `json:"moderator_id,omitempty"`
	CreatedAt          time.Time            `json:"created_at"`
	FormedAt           *time.Time           `json:"formed_at,omitempty"`
	ClosedAt           *time.Time           `json:"closed_at,omitempty"`
	Classes            []Class              `gorm:"many2many:mass_request_to_class;" json:"classes,omitempty"`
	MassRequestToClass []MassRequestToClass `gorm:"foreignKey:RequestID" json:"mass_request_to_class,omitempty"`
}