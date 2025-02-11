package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Kontentski/develops-today-task/internal/entity"
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/pkg/postgresql"
)

var _ service.TargetStorage = (*targetStorage)(nil)

type targetStorage struct {
	*postgresql.PostgreSQLGorm
}

func NewTargetStorage(postgresql *postgresql.PostgreSQLGorm) *targetStorage {
	return &targetStorage{postgresql}
}

func (s *targetStorage) GetTarget(ctx context.Context, id string) (*entity.Target, error) {
	var target entity.Target
	err := s.DB.First(&target, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get target: %w", err)
	}
	return &target, nil
}

func (s *targetStorage) CreateTarget(ctx context.Context, target *entity.Target) (*entity.Target, error) {
	err := s.DB.Create(target).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create target: %w", err)
	}
	return target, nil
}

func (s *targetStorage) UpdateTarget(ctx context.Context, target *entity.Target) (*entity.Target, error) {
	err := s.DB.Save(target).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update target: %w", err)
	}
	return target, nil
}

func (s *targetStorage) DeleteTarget(ctx context.Context, id string) error {
	err := s.DB.Where("id = ?", id).Delete(&entity.Target{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}
	return nil
}

func (s *targetStorage) ListTargets(ctx context.Context, missionID string) ([]entity.Target, error) {
	var targets []entity.Target
	err := s.DB.Where("mission_id = ?", missionID).Find(&targets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list targets: %w", err)
	}
	return targets, nil
}
