package rentals

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"tournaments/roles"
)

type Server struct {
	Db          *gorm.DB
	Router      *mux.Router
	AuthN       Authenticator
	AuthZ       *roles.Authorizer
	Initialized bool
}

func (s *Server) Setup() error {
	// Add all routes for resources
	resources := []Resource{
		&UserResource{Db: s.Db},
		&ApartmentResource{Db: s.Db},
	}
	for _, resource := range resources {
		CreateRoutes(resource, s.Router)
	}

	// Handlers that don't belong to resources
	s.Router.HandleFunc("/login", s.LoginHandler()).Methods("POST")

	// Add Authentication/Authorization middleware
	s.Router.Use(s.AuthenticationMiddleware)

	// Add content-type=application/json middleware
	s.Router.Use(s.ContentTypeJsonMiddleware)

	// Perform database migrations
	s.Db.AutoMigrate(DbModels...)

	// Initialize roles' permissions
	s.setupAuthorization()

	s.Initialized = true

	return nil
}

func (s *Server) setupAuthorization() {
	s.AuthZ.AddPermission("admin", "users", roles.Create, roles.Read, roles.Update, roles.Delete)
	s.AuthZ.AddPermission("admin", "apartments", roles.Create, roles.Read, roles.Update, roles.Delete)
	s.AuthZ.AddPermission("realtor", "apartments", roles.Create, roles.Read, roles.Update, roles.Delete)
	s.AuthZ.AddPermission("client", "apartments", roles.Read)
}
