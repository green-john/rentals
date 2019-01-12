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

	// Finds all resources of this type. query is
	// a querystring used to filter resources of this type.
	// In order to get all, use an empty query. Fields in query
	// correspond to JSON field names
	Find(query string) ([]byte, error)

	// Update the resource with the given jsonData.
	// Returns the updated resource or an error
	Update(id string, jsonData []byte) ([]byte, error)

	// Deletes the resource from the Db
	Delete(id string) error
}

// Creates all routes for the given resource.
// HTTP verbs are mapped to CRUD operations, respectively.
//
// POST -> Create
// GET -> Read
// PATCH -> Update
// DELETE -> Delete
func CreateRoutes(resource Resource, router *mux.Router) {
	getHandler := handleGet(resource)
	getAllHandler := handleGetAll(resource)
	postHandler := handlePost(resource)
	patchHandler := handlePatch(resource)
	deleteHandler := handleDelete(resource)

	url := fmt.Sprintf("/%s", resource.Name())
	urlWithId := fmt.Sprintf("%s/{id:[0-9]+}", url)

	router.HandleFunc(url, postHandler).Methods("POST")
	router.HandleFunc(url, getAllHandler).Methods("GET")
	router.HandleFunc(urlWithId, getHandler).Methods("GET")
	router.HandleFunc(urlWithId, patchHandler).Methods("PATCH")
	router.HandleFunc(urlWithId, deleteHandler).Methods("DELETE")
}

func handleGet(resource Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		result, err := resource.Read(vars["id"])
		if err != nil {
			badRequestError(err, w)
			return
		}

		_, err = w.Write([]byte(result))
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}
}

func handleGetAll(resource Resource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := resource.Find("")
		if err != nil {
			badRequestError(err, w)
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
			badRequestError(err, w)
			return
		}

		result, err := resource.Create(body)
		if err != nil {
			badRequestError(err, w)
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
			badRequestError(err, w)
			return
		}

		result, err := resource.Update(vars["id"], body)
		if err != nil {
			badRequestError(err, w)
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
			badRequestError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func badRequestError(err error, w http.ResponseWriter) {
	log.Printf("[ERROR] %s", err.Error())
	switch v := err.(type) {
	case NotFoundError:
		ErrorResponse(w, http.StatusNotFound, v.Error())
	default:
		ErrorResponse(w, http.StatusBadRequest, v.Error())
	}
}
