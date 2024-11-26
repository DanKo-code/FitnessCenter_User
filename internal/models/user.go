package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	Role        string    `db:"role"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
}
