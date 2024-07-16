package models

import "errors"

type Role struct {
	slug string
}

func (r Role) String() string {
	return r.slug
}

var (
	Unknown = Role{""}
	Member  = Role{"member"}
	Admin   = Role{"admin"}
)

func RoleFromString(s string) (Role, error) {
	switch s {
	case Member.slug:
		return Member, nil
	case Admin.slug:
		return Admin, nil
	default:
		return Unknown, errors.New("unknown role: " + s)
	}
}
