package entity

import (
	"time"

	"gorm.io/gorm"
)

// Mission represents a mission undertaken by a spy cat.
type Mission struct {
	ID        string         `json:"id,omitempty" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" binding:"required"`
	SpyCatID  *string        `json:"spyCatId,omitempty"`
	SpyCat    *SpyCat         `json:"spyCat,omitempty" gorm:"foreignKey:SpyCatID"`
	Targets   []Target       `json:"targets" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Completed bool           `json:"completed" binding:"required"`
	CreatedAt time.Time      `json:"createdAt,omitempty" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}
