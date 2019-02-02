package transport

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Server struct {
	Db      *gorm.DB
	handler http.Handler
	AuthN   AuthnService
	AuthZ   *AuthzService
}

// Creates an http server and serves it in the specified address
func (s *Server) ServeHTTP(addr string) error {
	srv := &http.Server{
		Handler:      s.handler,
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
		Db:      db,
		handler: router,
		AuthN:   authNService,
		AuthZ:   authZService,
	}

	for name, resource := range map[string]Resource{
		"users":      &UserResource{Db: db},
		"apartments": &ApartmentResource{Db: db},
	} {
		CreateRoutes(name, resource, router)
	}

	// Add other handlers
	router.HandleFunc("/login", s.LoginHandler()).Methods("POST")
	router.HandleFunc("/profile", s.profileHandler()).Methods("GET")
	router.HandleFunc("/newClient", s.newClientHandler()).Methods("POST")

	// Add Authentication/Authorization middleware
	router.Use(s.AuthMiddleware)

	// Add content-type=application/json middleware
	router.Use(s.ContentTypeJsonMiddleware)

	// Add CORS middleware
	router.Use(s.CORSMiddleware)

	// Log all things
	router.Use(s.LoggingMiddleware)

	// Initialize roles' permissions
	s.setupAuthorization()

	return s, nil
}
