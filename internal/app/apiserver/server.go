package apiserver

import (
	"github.com/antnzr/test-go-api/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	server := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	server.configureRouter()

	return server
}

func (server *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.router.ServeHTTP(writer, request)
}

func (server *server) configureRouter() {
	server.router.HandleFunc("/users", server.handleUsersCreate()).Methods("POST")
}

func (server *server) handleUsersCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
