package store

import "github.com/antnzr/test-go-api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(string) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
