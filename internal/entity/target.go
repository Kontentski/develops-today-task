package entity

import (
	"time"

	"gorm.io/gorm"
)
// Target represents a target within a mission.
type Target struct {
	ID        string         `json:"id,omitempty" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" binding:"required"`
	MissionID string         `json:"missionId" gorm:"type:uuid;not null"`
	Mission   *Mission       `json:"mission,omitempty" gorm:"foreignKey:MissionID"`
	Name      string         `json:"name" binding:"required"`
	Country   string         `json:"country" binding:"required"`
	Notes     string         `json:"notes"`
	Completed bool           `json:"completed" binding:"required"`
	CreatedAt time.Time      `json:"createdAt,omitempty" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}
