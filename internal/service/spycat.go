package service

import (
	"context"

	"github.com/Kontentski/develops-today-task/internal/entity"
	"gorm.io/gorm"
)

type spyCatService struct {
	serviceContext
	apis APIs
}

func NewSpyCatService(options Options, storage SpyCatStorage) SpyCatService {
	return &spyCatService{
		serviceContext: serviceContext{
			storages: options.Storages,
			cfg:      options.Config,
			apis:     options.APIs,
			logger:   options.Logger.Named("SpyCatService"),
		},
		apis: options.APIs,
	}
}

type CreateSpyCatOptions struct {
	Name              string
	YearsOfExperience int
	Breed             string
	Salary            float64
}

func (s *spyCatService) CreateSpyCat(ctx context.Context, opts CreateSpyCatOptions) (*entity.SpyCat, error) {
	s.logger.Info("Creating new spy cat", "opts", opts)

	// Validate breed via TheCatAPI
	if err := s.validateBreed(opts.Breed); err != nil {
		s.logger.Error("Invalid breed", "err", err)
		return nil, err
	}

	cat := &entity.SpyCat{
		Name:              opts.Name,
		YearsOfExperience: opts.YearsOfExperience,
		Breed:             opts.Breed,
		Salary:            opts.Salary,
	}

	createdCat, err := s.storages.SpyCat.CreateSpyCat(ctx, cat)
	if err != nil {
		s.logger.Error("Failed to create spy cat", "err", err)
		return nil, err
	}

	s.logger.Info("Spy cat created successfully", "cat", createdCat)
	return createdCat, nil
}

func (s *spyCatService) validateBreed(breed string) error {
	breeds, err := s.apis.CatAPI.GetBreeds()
	if err != nil {
		return err
	}

	for _, b := range breeds {
		if b.Name == breed {
			return nil
		}
	}
	return ErrCreateSpyCatInvalidBreed
}

func (s *spyCatService) DeleteSpyCat(ctx context.Context, id string) error {
	s.logger.Info("Deleting spy cat", "id", id)

	err := s.storages.SpyCat.DeleteSpyCat(ctx, id)
	if err == gorm.ErrRecordNotFound {
		s.logger.Error("Failed to delete spy cat", "err", err)
		return ErrDeleteSpyCatNotFound
	}

	s.logger.Info("Spy cat deleted successfully", "id", id)
	return nil
}

func (s *spyCatService) UpdateSpyCatSalary(ctx context.Context, id string, newSalary float64) (*entity.SpyCat, error) {
	s.logger.Info("Updating spy cat salary", "id", id, "newSalary", newSalary)

	cat, err := s.storages.SpyCat.GetSpyCat(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get spy cat", "err", err)
		return nil, err
	}
	if cat == nil {
		s.logger.Info("Spy cat not found", "id", id)
		return nil, ErrUpdateSpyCatNotFound
	}

	cat.Salary = newSalary
	updatedCat, err := s.storages.SpyCat.UpdateSpyCat(ctx, cat)
	if err != nil {
		s.logger.Error("Failed to update spy cat", "err", err)
		return nil, err
	}

	s.logger.Info("Spy cat salary updated successfully", "cat", updatedCat)
	return updatedCat, nil
}

func (s *spyCatService) ListSpyCats(ctx context.Context) ([]entity.SpyCat, error) {
	s.logger.Info("Listing all spy cats")

	cats, err := s.storages.SpyCat.ListSpyCats(ctx)
	if err != nil {
		s.logger.Error("Failed to list spy cats", "err", err)
		return nil, err
	}

	s.logger.Info("Spy cats listed successfully", "count", len(cats))
	return cats, nil
}

func (s *spyCatService) GetSpyCat(ctx context.Context, id string) (*entity.SpyCat, error) {
	s.logger.Info("Fetching spy cat", "id", id)

	cat, err := s.storages.SpyCat.GetSpyCat(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrGetSpyCatNotFound
		}
		s.logger.Error("Failed to get spy cat", "err", err)
		return nil, err
	}

	s.logger.Info("Spy cat fetched successfully", "cat", cat)
	return cat, nil
}
