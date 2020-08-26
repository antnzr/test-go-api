package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestStore(t *testing.T, databaseUrl string) (*Store, func(...string)) {
	t.Helper()
	config := NewConfig()

	config.DatabaseName = databaseUrl
	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if _, err := s.db.Collection("users").DeleteMany(ctx, bson.M{}); err != nil {
				t.Fatal(err)
			}
		}

		s.Close()
	}

}
