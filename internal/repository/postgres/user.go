package postgres

import (
	"User/internal/models"
	log_c "User/pkg/logger"
	"context"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRep *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := userRep.db.NamedExecContext(ctx, `
		UPDATE "user" 
		SET name = :name,
		    email = :email,
		    role = :role,
		    photo = :photo,
		    updated_time = :updated_time
		WHERE id = :id`, *user)
	if err != nil {
		log_c.ErrorLogger.Printf("Error UpdateUser: %v", err)
		return err
	}

	return nil
}
