package apiserver

import (
	"context"
	"github.com/antnzr/test-go-api/internal/app/store/mongostore"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

func Start(config *Config) error {
	db, err := newDb(config.DatabaseUrl)
	if err != nil {
		return err
	}

	store := mongostore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	server := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, server)
}

func newDb(databaseUrl string) (*mongo.Database, error) {
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

	db := client.Database("test-go-api")

	return db, nil
}

/*type ApiServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *mongostore.Store
}

func New(config *Config) *ApiServer {
	return &ApiServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *ApiServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("Starting api server...")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *ApiServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *ApiServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *ApiServer) handleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "Hello")
	}
}

func (s *ApiServer) configureStore() error {
	st := mongostore.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}
*/
