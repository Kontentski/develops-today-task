package service

import (
	"context"

	"github.com/Kontentski/develops-today-task/internal/entity"
	"github.com/Kontentski/develops-today-task/pkg/errs"
)

type missionService struct {
	serviceContext
	storage MissionStorage
}

func NewMissionService(options Options, storage MissionStorage) MissionService {
	return &missionService{
		serviceContext: serviceContext{
			storages: options.Storages,
			cfg:      options.Config,
			apis:     options.APIs,
			logger:   options.Logger.Named("MissionService"),
		},
		storage: storage,
	}
}

type CreateMissionOptions struct {
	Completed bool
	Targets   []CreateTargetOptions
}

func (s *missionService) CreateMission(ctx context.Context, opts CreateMissionOptions) (*entity.Mission, error) {
	s.logger.Info("Creating new mission", "opts", opts)

	if len(opts.Targets) == 0 || len(opts.Targets) > 3 {
		return nil, ErrCreateMissionInvalidTargets
	}

	mission := &entity.Mission{
		Completed: opts.Completed,
		Targets:   make([]entity.Target, len(opts.Targets)),
	}

	// Create targets
	for i, targetOpt := range opts.Targets {
		mission.Targets[i] = entity.Target{
			Name:      targetOpt.Name,
			Country:   targetOpt.Country,
			Notes:     targetOpt.Notes,
			Completed: targetOpt.Completed,
		}
	}

	createdMission, err := s.storage.CreateMission(ctx, mission)
	if err != nil {
		s.logger.Error("Failed to create mission", "err", err)
		return nil, err
	}

	s.logger.Info("Mission created successfully", "mission", createdMission)
	return createdMission, nil
}

func (s *missionService) DeleteMission(ctx context.Context, id string) error {
	s.logger.Info("Deleting mission", "id", id)

	mission, err := s.storage.GetMission(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get mission", "err", err)
		return err
	}
	if mission == nil {
		return ErrDeleteMissionNotFound
	}

	if mission.SpyCatID != nil {
		return ErrDeleteMissionAssigned
	}

	if err := s.storage.DeleteMission(ctx, id); err != nil {
		s.logger.Error("Failed to delete mission", "err", err)
		return err
	}

	s.logger.Info("Mission deleted successfully", "id", id)
	return nil
}

func (s *missionService) GetMission(ctx context.Context, id string) (*entity.Mission, error) {
	s.logger.Info("Fetching mission", "id", id)

	mission, err := s.storage.GetMission(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get mission", "err", err)
		return nil, err
	}
	if mission == nil {
		return nil, ErrGetMissionNotFound
	}

	s.logger.Info("Mission fetched successfully", "mission", mission)
	return mission, nil
}

type UpdateMissionOptions struct {
	Completed bool
}

func (s *missionService) UpdateMission(ctx context.Context, id string, opts UpdateMissionOptions) (*entity.Mission, error) {
	s.logger.Info("Updating mission", "id", id, "opts", opts)

	mission, err := s.storage.GetMission(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get mission", "err", err)
		return nil, err
	}
	if mission == nil {
		return nil, errs.New("mission not found")
	}

	mission.Completed = opts.Completed
	updatedMission, err := s.storage.UpdateMission(ctx, mission)
	if err != nil {
		s.logger.Error("Failed to update mission", "err", err)
		return nil, err
	}

	s.logger.Info("Mission updated successfully", "mission", updatedMission)
	return updatedMission, nil
}

func (s *missionService) ListMissions(ctx context.Context) ([]entity.Mission, error) {
	s.logger.Info("Listing all missions")

	missions, err := s.storage.ListMissions(ctx)
	if err != nil {
		s.logger.Error("Failed to list missions", "err", err)
		return nil, err
	}

	s.logger.Info("Missions listed successfully", "count", len(missions))
	return missions, nil
}

func (s *missionService) AssignSpyCat(ctx context.Context, missionID, spyCatID string) error {
	s.logger.Info("Assigning spy cat to mission", "missionID", missionID, "spyCatID", spyCatID)

	// Get mission
	mission, err := s.storage.GetMission(ctx, missionID)
	if err != nil {
		s.logger.Error("Failed to get mission", "err", err)
		return err
	}
	if mission == nil {
		return errs.New("mission not found")
	}

	if mission.SpyCatID != nil {
		return errs.New("mission already has an assigned cat")
	}

	// Get spy cat
	cat, err := s.storages.SpyCat.GetSpyCat(ctx, spyCatID)
	if err != nil {
		s.logger.Error("Failed to get spy cat", "err", err)
		return err
	}
	if cat == nil {
		return errs.New("spy cat not found")
	}

	if cat.MissionID != nil {
		return errs.New("spy cat is already assigned to a mission")
	}

	// Update mission with spy cat
	mission.SpyCatID = &spyCatID
	if _, err := s.storage.UpdateMission(ctx, mission); err != nil {
		s.logger.Error("Failed to update mission", "err", err)
		return err
	}

	s.logger.Info("Spy cat assigned to mission successfully")
	return nil
}
