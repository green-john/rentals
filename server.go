package rentals

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"tournaments/roles"
)

type Server struct {
	db          *gorm.DB
	router      *mux.Router
	authN       Authenticator
	authZ       *roles.Authorizer
	initialized bool
}

func (s *Server) Setup() error {
	// Add all routes for resources
	resources := []Resource{
		&UserResource{Db: s.db},
		&ApartmentResource{Db: s.db},
	}
	for _, resource := range resources {
		CreateRoutes(resource, s.router)
	}

	// Handlers that don't belong to resources
	s.router.HandleFunc("/login", s.LoginHandler())

	// Add authN middleware
	s.router.Use(s.AuthenticationMiddleware)

	// Perform database migrations
	s.db.AutoMigrate(DbModels...)

	// Initialize roles
	s.initRoles()

	s.initialized = true

	return nil
}

func (s *Server) initRoles() {
	s.authZ.AddPermission("admin", "users", roles.Create, roles.Read, roles.Update, roles.Delete)
	s.authZ.AddPermission("admin", "apartments", roles.Create, roles.Read, roles.Update, roles.Delete)
	s.authZ.AddPermission("realtor", "apartments", roles.Create, roles.Read, roles.Update, roles.Delete)
	s.authZ.AddPermission("client", "apartments", roles.Read)
}
