package user

import "errors"

var (
	ErrAdminRoleRequired  = errors.New("admin role is required")
	ErrUserRoleRequired   = errors.New("user role is required")
	ErrEmailAlreadyTaken  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound           = errors.New("user not found")
	ErrWrongPassword      = errors.New("wrong password")
	ErrEmptyPassHash      = errors.New("empty password hash")
)
