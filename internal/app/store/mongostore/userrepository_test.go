package mongostore_test

import (
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/antnzr/test-go-api/internal/app/store/mongostore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := mongostore.TestDB(t, databaseUrl, "test-go-api")
	defer teardown("users")
	s := mongostore.New(db)

	err := s.User().Create(&model.User{
		Email:    "ant@gmail.net",
		Password: "adminadm",
	})

	assert.NoError(t, err)
}
