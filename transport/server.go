package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"rentals"
	"rentals/services"
	"time"
)

// Implemented by models that only expose certain fields.
// This method returns a struct with json tags used by
// the transports
type Public interface {
	// The public method returns a struct
	// with json tags.
	Public() interface{}
}

type Server struct {
	Db               *gorm.DB
	router           *mux.Router
	authn            services.AuthnService
	authz            *services.AuthzService
	apartmentService services.ApartmentService
	userService      services.UserService
}

// Creates an http server and serves it in the specified address
func (s *Server) ServeHTTP(addr string) error {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) setupAuthorization() {
	s.authz.AddPermission("admin", "users",
		services.Create, services.Read, services.Update, services.Delete)
	s.authz.AddPermission("admin", "apartments",
		services.Create, services.Read, services.Update, services.Delete)
	s.authz.AddPermission("realtor", "apartments",
		services.Create, services.Read, services.Update, services.Delete)
	s.authz.AddPermission("client", "apartments", services.Read)
}

// Creates GET, POST, PATH and DELETE user handlers.
func (s *Server) AddApartmentsHandlers(basePath string) {
	url := fmt.Sprintf("/%s", basePath)
	urlWithId := fmt.Sprintf("%s/{id:[0-9]+}", url)

	s.router.HandleFunc(url, postApartmentsHandler(s.apartmentService)).Methods("POST")
	s.router.HandleFunc(url, getAllApartmentsHandler(s.apartmentService)).Methods("GET")
	s.router.HandleFunc(urlWithId, getApartmentsHandler(s.apartmentService)).Methods("GET")
	s.router.HandleFunc(urlWithId, patchApartmentsHandler(s.apartmentService)).Methods("PATCH")
	s.router.HandleFunc(urlWithId, deleteApartmentsHandler(s.apartmentService)).Methods("DELETE")
}

// Creates GET, POST, PATH and DELETE user handlers.
func (s *Server) AddUsersHandlers(basePath string) {
	url := fmt.Sprintf("/%s", basePath)
	urlWithId := fmt.Sprintf("%s/{id:[0-9]+}", url)

	s.router.HandleFunc(url, postUsersHandler(s.userService)).Methods("POST")
	s.router.HandleFunc(url, getAllUsersHandler(s.userService)).Methods("GET")
	s.router.HandleFunc(urlWithId, getUsersHandler(s.userService)).Methods("GET")
	s.router.HandleFunc(urlWithId, patchUsersHandler(s.userService)).Methods("PATCH")
	s.router.HandleFunc(urlWithId, deleteUsersHandler(s.userService)).Methods("DELETE")
}

func getUsersHandler(service services.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var input services.UserReadInput

		input.Id = vars["id"]
		result, err := service.Read(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func getAllUsersHandler(service services.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.UserAllInput

		result, err := service.All(input)
		if err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func postUsersHandler(service services.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var newUser services.UserCreateInput
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := service.Create(newUser)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusCreated, result)
	}
}

func patchUsersHandler(service services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()

		var updateInput services.UserUpdateInput

		if err := json.NewDecoder(r.Body).Decode(&updateInput); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			log.Println(err.Error())
			return
		}

		updateInput.Id = vars["id"]
		result, err := service.Update(updateInput)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func deleteUsersHandler(service services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var deleteIn services.UserDeleteInput
		deleteIn.Id = vars["id"]

		_, err := service.Delete(deleteIn)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusNoContent, nil)
	}
}

func getApartmentsHandler(srv services.ApartmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var input services.ApartmentReadInput
		input.Id = vars["id"]

		result, err := srv.Read(input)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func getAllApartmentsHandler(srv services.ApartmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.ApartmentFindInput
		input.Query = r.URL.RawQuery

		result, err := srv.Find(input)
		if err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func postApartmentsHandler(srv services.ApartmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var newApartment rentals.Apartment
		if err := json.NewDecoder(r.Body).Decode(&newApartment); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := srv.Create(services.ApartmentCreateInput{Apartment: newApartment})
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusCreated, result)
	}
}

func patchApartmentsHandler(srv services.ApartmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()

		var updateInput services.ApartmentUpdateInput

		if err := json.NewDecoder(r.Body).Decode(&updateInput.Data); err != nil {
			respond(w, http.StatusBadRequest, err.Error())
			return
		}

		updateInput.Id = vars["id"]

		result, err := srv.Update(updateInput)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusOK, result)
	}
}

func deleteApartmentsHandler(srv services.ApartmentService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteIn services.ApartmentDeleteInput

		vars := mux.Vars(r)
		deleteIn.Id = vars["id"]
		_, err := srv.Delete(deleteIn)
		if err != nil {
			badRequestError(err, w)
			return
		}

		respond(w, http.StatusNoContent, nil)
	}
}

func badRequestError(err error, w http.ResponseWriter) {
	log.Printf("[ERROR] %s", err.Error())
	switch err {
	case services.NotFoundError:
		respond(w, http.StatusNotFound, err.Error())
	default:
		respond(w, http.StatusBadRequest, err.Error())
	}
}

func NewServer(db *gorm.DB, authNService services.AuthnService, authZService *services.AuthzService,
	apartmentsService services.ApartmentService, userService services.UserService) (*Server, error) {
	router := mux.NewRouter()

	s := &Server{
		Db:               db,
		router:           router,
		authn:            authNService,
		authz:            authZService,
		apartmentService: apartmentsService,
		userService:      userService,
	}

	// Adds POST, GET, PATCH, DELETE for users
	s.AddUsersHandlers("users")

	// Adds POST, GET, PATCH, DELETE for apartments
	s.AddApartmentsHandlers("apartments")

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
