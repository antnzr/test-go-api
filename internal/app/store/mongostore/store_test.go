package mongostore_test

import (
	"os"
	"testing"
)

var (
	databaseUrl string
)

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "mongodb://localhost:27017/test-go-api-test"
	}

	os.Exit(m.Run())
}
