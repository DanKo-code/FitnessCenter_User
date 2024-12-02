package usecase

import (
	"User/internal/dtos"
	"User/internal/models"
	"context"
	"github.com/google/uuid"
)

type CoachUseCase interface {
	UpdateCoach(ctx context.Context, cmd *dtos.UpdateCoachCommand) (*models.Coach, error)
	CreateCoach(ctx context.Context, cmd *dtos.CreateCoachCommand) (*models.Coach, error)
	DeleteCoachById(ctx context.Context, id uuid.UUID) (*models.Coach, error)
	GetCoaches(ctx context.Context) ([]*models.Coach, error)
	GetCoachById(ctx context.Context, uuid uuid.UUID) (*models.Coach, error)
}
