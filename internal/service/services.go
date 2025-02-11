package service

import (
	"context"

	"github.com/Kontentski/develops-today-task/config"
	"github.com/Kontentski/develops-today-task/pkg/errs"
	"github.com/Kontentski/develops-today-task/pkg/logging"

	"github.com/Kontentski/develops-today-task/internal/entity"
)

type Services struct {
	SpyCat  SpyCatService
	Mission MissionService
	Target  TargetService
}

// serviceContext provides a shared context for all services
type serviceContext struct {
	storages Storages
	cfg      *config.Config
	logger   logging.Logger
	apis     APIs
}

// Options is used to parameterize service
type Options struct {
	Storages Storages
	APIs     APIs
	Config   *config.Config
	Logger   logging.Logger
}

// SpyCat errors
var (
	ErrCreateSpyCatInvalidBreed = errs.New("invalid breed")
	ErrDeleteSpyCatNotFound     = errs.New("spy cat not found")
	ErrUpdateSpyCatNotFound     = errs.New("spy cat not found")
	ErrGetSpyCatNotFound        = errs.New("spy cat not found")
)

// Mission errors
var (
	ErrCreateMissionInvalidTargets = errs.New("mission must have between 1 and 3 targets")
	ErrGetMissionNotFound          = errs.New("mission not found")
	ErrDeleteMissionNotFound       = errs.New("mission not found")
	ErrDeleteMissionAssigned       = errs.New("cannot delete mission assigned to a cat")
	ErrUpdateMissionNotFound       = errs.New("mission not found")
	ErrAssignMissionNotFound       = errs.New("mission not found")
	ErrAssignMissionHasCat         = errs.New("mission already has an assigned cat")
	ErrAssignSpyCatNotFound        = errs.New("spy cat not found")
	ErrAssignSpyCatBusy            = errs.New("spy cat is already assigned to a mission")
)

// Target errors
var (
	ErrCreateTargetMissionNotFound  = errs.New("mission not found")
	ErrCreateTargetCompletedMission = errs.New("cannot add target to completed mission")
	ErrCreateTargetTooMany          = errs.New("mission cannot have more than 3 targets")
	ErrGetTargetNotFound            = errs.New("target not found")
	ErrUpdateTargetNotFound         = errs.New("target not found")
	ErrUpdateTargetCompletedMission = errs.New("cannot update target in completed mission")
	ErrDeleteTargetNotFound         = errs.New("target not found")
	ErrDeleteTargetCompleted        = errs.New("cannot delete completed target")
)

// SpyCatService defines service operations for SpyCat.
type SpyCatService interface {
	CreateSpyCat(ctx context.Context, opts CreateSpyCatOptions) (*entity.SpyCat, error)
	GetSpyCat(ctx context.Context, id string) (*entity.SpyCat, error)
	UpdateSpyCatSalary(ctx context.Context, id string, newSalary float64) (*entity.SpyCat, error)
	DeleteSpyCat(ctx context.Context, id string) error
	ListSpyCats(ctx context.Context) ([]entity.SpyCat, error)
}

// MissionService defines service operations for Mission.
type MissionService interface {
	CreateMission(ctx context.Context, opts CreateMissionOptions) (*entity.Mission, error)
	GetMission(ctx context.Context, id string) (*entity.Mission, error)
	UpdateMission(ctx context.Context, id string, opts UpdateMissionOptions) (*entity.Mission, error)
	DeleteMission(ctx context.Context, id string) error
	ListMissions(ctx context.Context) ([]entity.Mission, error)
	AssignSpyCat(ctx context.Context, missionID, spyCatID string) error
}

// TargetService defines service operations for Target.
type TargetService interface {
	CreateTarget(ctx context.Context, missionID string, opts CreateTargetOptions) (*entity.Target, error)
	UpdateTarget(ctx context.Context, id string, opts UpdateTargetOptions) (*entity.Target, error)
	DeleteTarget(ctx context.Context, id string) error
	ListTargets(ctx context.Context, missionID string) ([]entity.Target, error)
}

func NewService(options Options) Services {
	return Services{
		SpyCat:  NewSpyCatService(options, options.Storages.SpyCat),
		Mission: NewMissionService(options, options.Storages.Mission),
		Target:  NewTargetService(options, options.Storages.Target),
	}
}
