package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Store struct {
	config         *Config
	db             *mongo.Database
	userRepository *UserRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	//clientOptions := options.Client().ApplyURI("mongodb://user:password@host:port/test?authSource=admin&replicaSet=Cluster0-shard-0&readPreference=primary&test-go-api=MongoDB%20Compass&ssl=true")
	clientOptions := options.Client().ApplyURI(s.config.DatabaseUrl)
	client, err := mongo.NewClient(clientOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	s.db = client.Database(s.config.DatabaseName)

	return nil
}

func (s *Store) Close() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.db.Client().Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}

func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
