package service

import (
	"context"

	"github.com/Kontentski/develops-today-task/internal/entity"
)

type targetService struct {
	serviceContext
	storage TargetStorage
}

func NewTargetService(options Options, storage TargetStorage) TargetService {
	return &targetService{
		serviceContext: serviceContext{
			storages: options.Storages,
			cfg:      options.Config,
			apis:     options.APIs,
			logger:   options.Logger.Named("TargetService"),
		},
		storage: storage,
	}
}

type CreateTargetOptions struct {
	Name      string
	Country   string
	Notes     string
	Completed bool
}

func (s *targetService) CreateTarget(ctx context.Context, missionID string, opts CreateTargetOptions) (*entity.Target, error) {
	s.logger.Info("Creating new target", "missionID", missionID, "opts", opts)

	mission, err := s.storages.Mission.GetMission(ctx, missionID)
	if err != nil {
		s.logger.Error("Failed to get mission", "err", err)
		return nil, err
	}
	if mission == nil {
		return nil, ErrCreateTargetMissionNotFound
	}

	if mission.Completed {
		return nil, ErrCreateTargetCompletedMission
	}

	if len(mission.Targets) >= 3 {
		return nil, ErrCreateTargetTooMany
	}

	target := &entity.Target{
		MissionID: missionID,
		Name:      opts.Name,
		Country:   opts.Country,
		Notes:     opts.Notes,
		Completed: opts.Completed,
	}

	createdTarget, err := s.storage.CreateTarget(ctx, target)
	if err != nil {
		s.logger.Error("Failed to create target", "err", err)
		return nil, err
	}

	s.logger.Info("Target created successfully", "target", createdTarget)
	return createdTarget, nil
}

type UpdateTargetOptions struct {
	Notes     *string
	Completed *bool
}

func (s *targetService) UpdateTarget(ctx context.Context, id string, opts UpdateTargetOptions) (*entity.Target, error) {
	s.logger.Info("Updating target", "id", id, "opts", opts)

	target, err := s.storage.GetTarget(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get target", "err", err)
		return nil, err
	}
	if target == nil {
		return nil, ErrUpdateTargetNotFound
	}

	mission, err := s.storages.Mission.GetMission(ctx, target.MissionID)
	if err != nil {
		s.logger.Error("Failed to get mission", "err", err)
		return nil, err
	}

	if opts.Notes != nil {
		if target.Completed {
			return nil, ErrUpdateTargetCompletedMission
		}
		if mission.Completed {
			return nil, ErrUpdateTargetCompletedMission
		}
		target.Notes = *opts.Notes
	}

	// Handle completion update
	if opts.Completed != nil {
		target.Completed = *opts.Completed

		// If all targets are completed, mark mission as completed
		if *opts.Completed {
			allCompleted := true
			for _, t := range mission.Targets {
				// check if all targets are completed except the current one
				if t.ID != target.ID && !t.Completed {
					allCompleted = false
					break
				}
			}
			if allCompleted {
				mission.Completed = true
				if _, err := s.storages.Mission.UpdateMission(ctx, mission); err != nil {
					s.logger.Error("Failed to update mission completion status", "err", err)
					return nil, err
				}
			}
		}
	}

	updatedTarget, err := s.storage.UpdateTarget(ctx, target)
	if err != nil {
		s.logger.Error("Failed to update target", "err", err)
		return nil, err
	}

	s.logger.Info("Target updated successfully", "target", updatedTarget)
	return updatedTarget, nil
}

func (s *targetService) DeleteTarget(ctx context.Context, id string) error {
	s.logger.Info("Deleting target", "id", id)

	target, err := s.storage.GetTarget(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get target", "err", err)
		return err
	}
	if target == nil {
		return ErrDeleteTargetNotFound
	}

	if target.Completed {
		return ErrDeleteTargetCompleted
	}

	if err := s.storage.DeleteTarget(ctx, id); err != nil {
		s.logger.Error("Failed to delete target", "err", err)
		return err
	}

	s.logger.Info("Target deleted successfully", "id", id)
	return nil
}

func (s *targetService) ListTargets(ctx context.Context, missionID string) ([]entity.Target, error) {
	s.logger.Info("Listing targets for mission", "missionID", missionID)

	targets, err := s.storage.ListTargets(ctx, missionID)
	if err != nil {
		s.logger.Error("Failed to list targets", "err", err)
		return nil, err
	}

	s.logger.Info("Targets listed successfully", "count", len(targets))
	return targets, nil
}
