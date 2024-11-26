package repository

import (
	"User/internal/models"
	"context"
)

type UserRepository interface {
	UpdateUser(ctx context.Context, user *models.User) error
}
