package service

import (
	"context"

	"github.com/Kontentski/develops-today-task/internal/entity"
)

type Storages struct {
	SpyCat  SpyCatStorage
	Mission MissionStorage
	Target  TargetStorage
}

// SpyCatStorage defines storage operations for SpyCat.
type SpyCatStorage interface {
	GetSpyCat(ctx context.Context, id string) (*entity.SpyCat, error)
	CreateSpyCat(ctx context.Context, cat *entity.SpyCat) (*entity.SpyCat, error)
	UpdateSpyCat(ctx context.Context, cat *entity.SpyCat) (*entity.SpyCat, error)
	DeleteSpyCat(ctx context.Context, id string) error
	ListSpyCats(ctx context.Context) ([]entity.SpyCat, error)
}

// MissionStorage defines storage operations for Mission.
type MissionStorage interface {
	GetMission(ctx context.Context, id string) (*entity.Mission, error)
	CreateMission(ctx context.Context, mission *entity.Mission) (*entity.Mission, error)
	UpdateMission(ctx context.Context, mission *entity.Mission) (*entity.Mission, error)
	DeleteMission(ctx context.Context, id string) error
	ListMissions(ctx context.Context) ([]entity.Mission, error)
}

// TargetStorage defines storage operations for Target.
type TargetStorage interface {
	GetTarget(ctx context.Context, id string) (*entity.Target, error)
	CreateTarget(ctx context.Context, target *entity.Target) (*entity.Target, error)
	UpdateTarget(ctx context.Context, target *entity.Target) (*entity.Target, error)
	DeleteTarget(ctx context.Context, id string) error
	ListTargets(ctx context.Context, missionID string) ([]entity.Target, error)
}
