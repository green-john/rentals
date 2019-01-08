package rentals

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	db          *gorm.DB
	router      *mux.Router
	authN       Authenticator
	initialized bool
}

func (s *Server) Setup() error {
	// Add all routes for resources
	resources := []Resource{&UserResource{s.db}}
	for _, resource := range resources {
		CreateRoutes(resource, s.router)
	}

	// Handlers that don't belong to resources
	s.router.HandleFunc("/login", s.LoginHandler())

	// Add authN middleware
	s.router.Use(s.AuthorizationMiddleware)

	// Perform database migrations
	s.db.AutoMigrate(DbModels...)

	s.initialized = true

	return nil
}
