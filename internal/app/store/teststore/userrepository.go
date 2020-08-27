package teststore

import (
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/antnzr/test-go-api/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (repo *UserRepository) Create(usr *model.User) error {
	if err := usr.Validate(); err != nil {
		return err
	}

	if err := usr.BeforeCreate(); err != nil {
		return err
	}

	repo.users[usr.Email] = usr
	usr.ID = string(len(repo.users))

	return nil
}

func (repo *UserRepository) FindByEmail(email string) (*model.User, error) {
	usr, ok := repo.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return usr, nil
}
