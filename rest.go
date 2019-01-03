package rentals

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// Represents a high level rest resource.
//
// See `CreateRoutes`
type Resource interface {
	// Name of the resource. To be used in the URL
	Name() string

	// Create a new resource given the json data. Returns
	// the created object or an error
	Create(jsonData []byte) ([]byte, error)

	// Get resource for the requested ID
	Read(id string) ([]byte, error)

	// Update the resource with the given jsonData.
	// Returns the updated resource or an error
	Update(id string, jsonData []byte) ([]byte, error)

	// Deletes the resource from the DB
	Delete(id string) error
}

// A wrapper around our http router so that it is easier to test
type HttpRouter interface {
	// Adds the url with the given method to the Router
	Add(url string, method string, fn func(http.ResponseWriter, *http.Request))

	// Use the given middleware in the router
	Use(middleware func(http.Handler) http.Handler)
}

type GorillaRouter struct {
	Router *mux.Router
}

func (r *GorillaRouter) Add(url string, method string, fn func(http.ResponseWriter, *http.Request)) {
	r.Router.HandleFunc(url, fn).Methods(method)
}

func (r *GorillaRouter) Use(middleware func(http.Handler) http.Handler) {
	r.Router.Use(middleware)
}

func NewGorillaRouter() *GorillaRouter {
	router := mux.NewRouter()
	return &GorillaRouter{
		Router: router,
	}
}

// Creates all routes for the given resource.
// HTTP verbs are mapped to CRUD operations, respectively.
//
// POST -> Create
// GET -> Read
// PATCH -> Update
// DELETE -> Delete
func CreateRoutes(resource Resource, router HttpRouter) {
	getHandler := handleGet(resource)
	postHandler := handlePost(resource)
	patchHandler := handlePatch(resource)
	deleteHandler := handleDelete(resource)

	url := fmt.Sprintf("/%s", resource.Name())
	urlWithId := fmt.Sprintf("%s/{id:[0-9]+}", url)

	router.Add(url, "POST", postHandler)
	router.Add(urlWithId, "GET", getHandler)
	router.Add(urlWithId, "PATCH", patchHandler)
	router.Add(urlWithId, "DELETE", deleteHandler)
}

func handleGet(resource Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		result, err := resource.Read(vars["id"])
		if err != nil {
			writeError(err, w)
			return
		}

		_, err = w.Write([]byte(result))
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}
}

func handlePost(resource Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(err, w)
			return
		}

		result, err := resource.Create(body)
		if err != nil {
			writeError(err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(result))
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}
}

func handlePatch(resource Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(err, w)
			return
		}

		result, err := resource.Update(vars["id"], body)
		if err != nil {
			writeError(err, w)
			return
		}

		_, err = w.Write([]byte(result))
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}
}

func handleDelete(resource Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := resource.Delete(vars["id"])
		if err != nil {
			writeError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func writeError(err error, w http.ResponseWriter) {
	w.WriteHeader(400)
	_, _ = w.Write([]byte(err.Error()))
	log.Printf("[ERROR] %s", err.Error())
}
