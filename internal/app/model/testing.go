package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Name:     "Nazar",
		Password: "adminadmin",
	}
}
