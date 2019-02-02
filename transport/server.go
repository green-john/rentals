package transport

import (
	"errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"rentals"
	"time"
)

type Server struct {
	Db          *gorm.DB
	Router      *mux.Router
	AuthN       AuthnService
	AuthZ       *AuthzService
	Initialized bool
}

func (s *Server) Setup() error {
	// Create routes for all resources
	for _, resource := range []Resource{
		&UserResource{Db: s.Db},
		&ApartmentResource{Db: s.Db},
	} {
		CreateRoutes(resource, s.Router)
	}

	// Add other handlers
	s.Router.HandleFunc("/login", s.LoginHandler()).Methods("POST")
	s.Router.HandleFunc("/profile", s.profileHandler()).Methods("GET")
	s.Router.HandleFunc("/newClient", s.newClientHandler()).Methods("POST")

	// Add Authentication/Authorization middleware
	s.Router.Use(s.AuthMiddleware)

	// Add content-type=application/json middleware
	s.Router.Use(s.ContentTypeJsonMiddleware)

	// Perform database migrations
	s.Db.AutoMigrate(rentals.DbModels...)

	// Initialize roles' permissions
	s.setupAuthorization()

	s.Initialized = true

	return nil
}

// Serves transport requests. app.Server must be Initialized,
// otherwise an error is thrown
func (s *Server) ServeHTTP(addr string) error {
	if !s.Initialized {
		return errors.New("app must be Initialized first")
	}

	// Enable CORS for testing purposes. This should be
	// configured properly for production
	allOrigins := handlers.AllowedOrigins([]string{"ruizandr.es,localhost"})
	allMethods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"})
	allHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	// Log all requests
	r := handlers.LoggingHandler(os.Stderr, s.Router)

	srv := &http.Server{
		Handler:      handlers.CORS(allOrigins, allMethods, allHeaders)(r),
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) setupAuthorization() {
	s.AuthZ.AddPermission("admin", "users", Create, Read, Update, Delete)
	s.AuthZ.AddPermission("admin", "apartments", Create, Read, Update, Delete)
	s.AuthZ.AddPermission("realtor", "apartments", Create, Read, Update, Delete)
	s.AuthZ.AddPermission("client", "apartments", Read)
}

func NewServer(db *gorm.DB, authNService AuthnService, authZService *AuthzService) (*Server, error) {
	router := mux.NewRouter()

	s := &Server{
		Db:          db,
		Router:      router,
		AuthN:       authNService,
		AuthZ:       authZService,
		Initialized: false,
	}

	err := s.Setup()
	if err != nil {
		return nil, err
	}

	return s, nil
}
