package store

import "github.com/antnzr/test-go-api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
}
