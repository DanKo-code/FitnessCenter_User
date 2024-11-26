package user_usecase

import (
	"User/internal/models"
	"User/internal/repository"
	"context"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) UpdateUser(ctx context.Context, uuReq *models.User) (*models.User, error) {

	err := u.userRepo.UpdateUser(ctx, uuReq)
	if err != nil {
		return nil, err
	}

	return uuReq, nil
}
