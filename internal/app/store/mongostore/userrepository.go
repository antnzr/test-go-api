package mongostore

import (
	"context"
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/antnzr/test-go-api/internal/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	data := bson.M{
		"name":          usr.Name,
		"email":         usr.Email,
		"password_hash": usr.PasswordHash,
	}

	_, err := collection.InsertOne(ctx, data)
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

func (r *UserRepository) Find(id string) (*model.User, error) {
	usr := &model.User{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	filter := bson.M{"_id": objectId}

	collection := r.store.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := collection.FindOne(ctx, filter).Decode(&usr); err != nil {
		return nil, store.ErrRecordNotFound
	}

	return usr, nil
}
