package util

import "github.com/WilSimpson/gocloak/v13"

func RegisterRole(role *gocloak.Role, roles *[]*gocloak.Role) *gocloak.Role {
	*roles = append(*roles, role)
	return role
}

