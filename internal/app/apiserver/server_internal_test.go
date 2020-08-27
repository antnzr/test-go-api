package apiserver

import (
	"github.com/antnzr/test-go-api/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)

	server := newServer(teststore.New())
	server.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
}
