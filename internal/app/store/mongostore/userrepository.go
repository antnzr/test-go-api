package mongostore

import (
	"context"
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/antnzr/test-go-api/internal/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(usr *model.User) error {
	if err := usr.Validate(); err != nil {
		return err
	}

	if err := usr.BeforeCreate(); err != nil {
		return err
	}

	collection := r.store.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := collection.InsertOne(ctx, bson.M{"name": usr.Name, "password_hash": usr.PasswordHash})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	usr := &model.User{}
	filter := bson.M{"email": email}

	collection := r.store.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := collection.FindOne(ctx, filter).Decode(&usr); err != nil {
		return nil, store.ErrRecordNotFound
	}

	return usr, nil
}
