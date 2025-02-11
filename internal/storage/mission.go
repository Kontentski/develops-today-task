package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Kontentski/develops-today-task/internal/entity"
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/pkg/postgresql"
)

var _ service.MissionStorage = (*missionStorage)(nil)

type missionStorage struct {
	*postgresql.PostgreSQLGorm
}

func NewMissionStorage(postgresql *postgresql.PostgreSQLGorm) *missionStorage {
	return &missionStorage{postgresql}
}

func (s *missionStorage) CreateMission(ctx context.Context, mission *entity.Mission) (*entity.Mission, error) {
	err := s.DB.Create(mission).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create mission: %w", err)
	}
	return mission, nil
}

func (s *missionStorage) GetMission(ctx context.Context, id string) (*entity.Mission, error) {
	var mission entity.Mission
	err := s.DB.
		Preload("SpyCat").
		Preload("Targets").
		First(&mission, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get mission: %w", err)
	}
	return &mission, nil
}

func (s *missionStorage) UpdateMission(ctx context.Context, mission *entity.Mission) (*entity.Mission, error) {
	err := s.DB.Save(mission).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update mission: %w", err)
	}
	return mission, nil
}

func (s *missionStorage) DeleteMission(ctx context.Context, id string) error {
	err := s.DB.Where("id = ?", id).Delete(&entity.Mission{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete mission: %w", err)
	}
	return nil
}

func (s *missionStorage) ListMissions(ctx context.Context) ([]entity.Mission, error) {
	var missions []entity.Mission
	err := s.DB.
		Preload("SpyCat").
		Preload("Targets").
		Find(&missions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list missions: %w", err)
	}
	return missions, nil
}
