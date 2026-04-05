package user

import "errors"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

var ErrInvalidRole = errors.New("invalid user role")

func ParseRole(v string) (Role, error) {
	r := Role(v)
	switch r {
	case RoleAdmin, RoleUser:
		return r, nil
	default:
		return "", ErrInvalidRole
	}
}
