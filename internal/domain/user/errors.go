package user

import "errors"

var (
	ErrAdminRoleRequired = errors.New("admin role is required")
)
