package postgres

import (
	"User/internal/models"
	"User/internal/repository"
	"User/pkg/logger"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var _ repository.CoachRepository = (*CoachRepository)(nil)

type CoachRepository struct {
	db *sqlx.DB
}

func NewCoachRepository(db *sqlx.DB) *CoachRepository {
	return &CoachRepository{db: db}
}

func (coachRep *CoachRepository) CreateCoach(ctx context.Context, coach *models.Coach) (*models.Coach, error) {
	_, err := coachRep.db.NamedExecContext(ctx, `
	INSERT INTO "coach" (id, name, description, photo, created_time, updated_time)
	VALUES (:id, :name, :description, :photo, :created_time, :updated_time)`, *coach)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateCoach: %v", err)
		return nil, err
	}

	return coach, nil
}

func (coachRep *CoachRepository) UpdateCoach(ctx context.Context, coach *models.Coach) error {
	_, err := coachRep.db.NamedExecContext(ctx, `
		UPDATE "coach"
		SET id = :id,
		    name = :name,
		    description = :description,
		    photo = :photo,
		    created_time = :created_time,
		    updated_time = :updated_time
		WHERE id = :id`, *coach)
	if err != nil {
		logger.ErrorLogger.Printf("Error UpdateUser: %v", err)
		return err
	}

	return nil
}

func (coachRep *CoachRepository) DeleteCoachById(ctx context.Context, id uuid.UUID) error {
	_, err := coachRep.db.NamedExecContext(ctx, `
		DELETE FROM "coach"
		WHERE id = :id`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteCoach: %v", err)
		return err
	}

	return nil
}

func (coachRep *CoachRepository) GetCoaches(ctx context.Context) ([]*models.Coach, error) {
	var coaches []*models.Coach

	err := coachRep.db.SelectContext(ctx, &coaches, `SELECT id, name, description, photo, created_time, updated_time FROM "coach"`)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetCoaches: %v", err)
		return nil, err
	}

	return coaches, nil
}

func (coachRep *CoachRepository) GetCoachById(ctx context.Context, id uuid.UUID) (*models.Coach, error) {
	var coach *models.Coach
	err := coachRep.db.GetContext(ctx, &coach, `
		SELECT id, name, description, photo, created_time, updated_time
		FROM "coach"
		WHERE id = :id`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetCoachById: %v", err)
		return nil, err
	}

	return coach, nil
}
