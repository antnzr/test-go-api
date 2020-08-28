package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/antnzr/test-go-api/internal/app/model"
	"github.com/antnzr/test-go-api/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

const (
	sessionName        = "superSessionName"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

func newServer(store store.Store, sessionStore sessions.Store) *server {
	server := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	server.configureRouter()

	return server
}

func (server *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.router.ServeHTTP(writer, request)
}

func (server *server) configureRouter() {
	server.router.Use(server.setRequestID)
	server.router.Use(server.logRequest)
	server.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	server.router.HandleFunc("/users", server.handleUsersCreate()).Methods("POST")
	server.router.HandleFunc("/sessions", server.handleSessionCreate()).Methods("POST")

	private := server.router.PathPrefix("/private").Subrouter()
	private.Use(server.authenticateUser)
	private.HandleFunc("/whoami", server.handleWhoami()).Methods("GET")
}

func (server *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (server *server) handleSessionCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(writer, r, http.StatusBadRequest, err)
		}

		usr, err := server.store.User().FindByEmail(req.Email)

		if err != nil || !usr.ComparePassword(req.Password) {
			server.error(writer, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := server.sessionStore.Get(r, sessionName)
		if err != nil {
			server.error(writer, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = usr.ID
		if err := server.sessionStore.Save(r, writer, session); err != nil {
			server.error(writer, r, http.StatusInternalServerError, err)
			return
		}

		server.respond(writer, r, http.StatusOK, nil)
	}
}

func (server *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (server *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := server.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (server *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := server.sessionStore.Get(r, sessionName)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
		}

		id, ok := session.Values["user_id"]
		if !ok {
			server.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		usr, err := server.store.User().Find(id.(string))

		if err != nil {
			server.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, usr)))
	})
}

func (server *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(writer, r, http.StatusBadRequest, err)
			return
		}

		usr := &model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := server.store.User().Create(usr); err != nil {
			server.error(writer, r, http.StatusUnprocessableEntity, err)
		}

		usr.Sanitize()
		server.respond(writer, r, http.StatusCreated, usr)
	}
}

func (server *server) error(writer http.ResponseWriter, request *http.Request, code int, err error) {
	server.respond(writer, request, code, map[string]string{"error": err.Error()})
}

func (server *server) respond(writer http.ResponseWriter, request *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)

	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}
