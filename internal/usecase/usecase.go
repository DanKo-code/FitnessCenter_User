package usecase

import (
	"User/internal/models"
	"context"
)

type UserUseCase interface {
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type CloudUseCase interface {
	PutObject(ctx context.Context, object []byte, name string) (string, error)
}
