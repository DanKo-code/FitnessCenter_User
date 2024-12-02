package errors

import "errors"

var (
	VoidUserData      = errors.New("void user data")
	VoidCoachData     = errors.New("void coach data")
	UserAlreadyExists = errors.New("user already exists")
	UserNotFound      = errors.New("user not found")
	InvalidPassword   = errors.New("invalid password")
)
