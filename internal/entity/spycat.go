package entity

import (
	"time"

	"gorm.io/gorm"
)

// SpyCat represents a spy cat in the system.
type SpyCat struct {
	ID                string         `json:"id,omitempty" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" binding:"required"`
	Name              string         `json:"name" binding:"required"`
	YearsOfExperience int            `json:"yearsOfExperience" binding:"required,gt=0"`
	Breed             string         `json:"breed" binding:"required"`
	Salary            float64        `json:"salary" binding:"required,gt=0"`
	MissionID         *string        `json:"missionId,omitempty" gorm:"type:uuid"`
	Mission           *Mission       `json:"mission,omitempty"`
	CreatedAt         time.Time      `json:"createdAt,omitempty" gorm:"index"`
	UpdatedAt         time.Time      `json:"updatedAt,omitempty"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
} 
