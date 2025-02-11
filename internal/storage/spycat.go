package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Kontentski/develops-today-task/internal/entity"
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/pkg/postgresql"
)

var _ service.SpyCatStorage = (*spyCatStorage)(nil)

type spyCatStorage struct {
	*postgresql.PostgreSQLGorm
}

func NewSpyCatStorage(postgresql *postgresql.PostgreSQLGorm) *spyCatStorage {
	return &spyCatStorage{postgresql}
}

func (s *spyCatStorage) CreateSpyCat(ctx context.Context, cat *entity.SpyCat) (*entity.SpyCat, error) {
	err := s.DB.Create(cat).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create spy cat: %w", err)
	}
	return cat, nil
}

func (s *spyCatStorage) UpdateSpyCat(ctx context.Context, cat *entity.SpyCat) (*entity.SpyCat, error) {
	err := s.DB.Save(cat).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update spy cat: %w", err)
	}
	return cat, nil
}

func (s *spyCatStorage) DeleteSpyCat(ctx context.Context, id string) error {
	err := s.DB.Where("id = ?", id).Delete(&entity.SpyCat{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete spy cat: %w", err)
	}
	return nil
}

func (s *spyCatStorage) ListSpyCats(ctx context.Context) ([]entity.SpyCat, error) {
	var cats []entity.SpyCat
	err := s.DB.Find(&cats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list spy cats: %w", err)
	}
	return cats, nil
}

func (s *spyCatStorage) GetSpyCat(ctx context.Context, id string) (*entity.SpyCat, error) {
	var cat entity.SpyCat
	err := s.DB.First(&cat, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get spy cat: %w", err)
	}
	return &cat, nil
}
