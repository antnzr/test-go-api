package store_test

import (
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/antnzr/test-go-api/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, _ := store.TestStore(t, databaseUrl)
	u, err := s.User().Create(&model.User{
		Name: "Antonio",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}
