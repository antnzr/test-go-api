package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Password     string
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (usr *User) Validate() error {
	return validation.ValidateStruct(
		usr,
		validation.Field(
			&usr.Password,
			validation.By(requiredIf(usr.PasswordHash == "")),
			validation.Length(6, 25),
		),
	)
}

func (usr *User) BeforeCreate() error {
	if len(usr.Password) > 0 {
		enc, err := encryptString(usr.Password)
		if err != nil {
			return err
		}

		usr.PasswordHash = enc
	}
	return nil
}

func encryptString(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
