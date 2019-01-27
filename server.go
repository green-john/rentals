package rentals

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"rentals/roles"
)

type Server struct {
	Db          *gorm.DB
	Router      *mux.Router
	AuthN       Authenticator
	AuthZ       *roles.Authorizer
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
