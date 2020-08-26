package model_test

import (
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	usr := model.TestUser(t)

	assert.NoError(t, usr.Validate())
}

func TestUser_BeforeCreate(t *testing.T) {
	usr := model.TestUser(t)

	assert.NoError(t, usr.BeforeCreate())
	logrus.Info(usr)
	assert.NotEmpty(t, usr.PasswordHash)
}
