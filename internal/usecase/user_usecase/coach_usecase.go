package user_usecase

import (
	"User/internal/dtos"
	"User/internal/models"
	"User/internal/repository"
	"context"
	"github.com/google/uuid"
	"time"
)

type CoachUseCase struct {
	coachRepo repository.CoachRepository
}

func NewCoachUseCase(coachRepo repository.CoachRepository) *CoachUseCase {
	return &CoachUseCase{coachRepo: coachRepo}
}

func (u *CoachUseCase) UpdateCoach(ctx context.Context, cmd *dtos.UpdateCoachCommand) (*models.Coach, error) {

	previousCoach, id := u.coachRepo.GetCoachById(ctx, cmd.Id)
	if id != nil {
		return nil, id
	}

	newCoach := &models.Coach{
		Id:          uuid.New(),
		Name:        cmd.Name,
		Description: cmd.Description,
		Photo:       cmd.Photo,
		UpdatedTime: time.Now(),
	}

	err := u.coachRepo.UpdateCoach(ctx, newCoach)
	if err != nil {
		return nil, err
	}

	updatedCoach := &models.Coach{
		Id:          previousCoach.Id,
		Name:        newCoach.Name,
		Description: newCoach.Description,
		Photo:       newCoach.Photo,
		UpdatedTime: newCoach.UpdatedTime,
		CreatedTime: previousCoach.CreatedTime,
	}

	return updatedCoach, nil
}
