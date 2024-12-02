package dtos

type CreateCoachCommand struct {
	Name        string `db:"name"`
	Description string `db:"description"`
}
