package mongostore

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T, databaseUrl string, databaseName string) (*mongo.Database, func(...string)) {
	t.Helper()
	//config := NewConfig()
	//
	//config.DatabaseName = databaseUrl
	//s := New(config)
	//if err := s.Open(); err != nil {
	//	t.Fatal(err)
	//}

	clientOptions := options.Client().ApplyURI(databaseUrl)
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

	db := client.Database(databaseName)

	return db, func(collections ...string) {
		if len(collections) > 0 {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if _, err := db.Collection("users").DeleteMany(ctx, bson.M{}); err != nil {
				t.Fatal(err)
			}
		}
	}

}
