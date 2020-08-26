package store

import (
	"context"
	"github.com/antnzr/test-go-api/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(usr *model.User) (*model.User, error) {
	if err := usr.Validate(); err != nil {
		return nil, err
	}

	if err := usr.BeforeCreate(); err != nil {
		return nil, err
	}

	collection := r.store.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := collection.InsertOne(ctx, bson.M{"name": usr.Name, "password_hash": usr.PasswordHash})
	if err != nil {
		return nil, err
	}

	usr.ID = res.InsertedID.(primitive.ObjectID).String()

	return usr, nil
}

func (r *UserRepository) FindByName(name string) (*model.User, error) {
	usr := &model.User{}
	filter := bson.M{"name": name}

	collection := r.store.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := collection.FindOne(ctx, filter).Decode(&usr); err != nil {
		return nil, err
	}

	return usr, nil
}
