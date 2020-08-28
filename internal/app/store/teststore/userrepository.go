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

	usr.ID = string(len(repo.users) + 1)
	repo.users[usr.ID] = usr

	return nil
}

func (repo *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, usr := range repo.users {
		if usr.Email == email {
			return usr, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (repo *UserRepository) Find(id string) (*model.User, error) {
	usr, ok := repo.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return usr, nil
}
