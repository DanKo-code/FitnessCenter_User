package usecase

import (
	"User/internal/dtos"
	"User/internal/models"
	"context"
	"github.com/google/uuid"
)

type UserUseCase interface {
	UpdateUser(ctx context.Context, cmd *dtos.UpdateUserCommand) (*models.User, error)
	CreateUser(ctx context.Context, cmd *dtos.CreateUserCommand) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	CheckPassword(ctx context.Context, cmd *dtos.CheckPasswordCommand) error
	DeleteUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
}
