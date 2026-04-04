package user

import "errors"

var (
	ErrAdminRoleRequired = errors.New("admin role is required")
	ErrUserRoleRequired  = errors.New("user role is required")
)
