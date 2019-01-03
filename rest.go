package rentals

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// Resource to be managed by the application. Methods map to
// HTTP verbs to comply with the REST standard
// POST -> New
// GET -> Fetch
// PATCH -> Update
// DELETE -> Remove
type Resource interface {
	// Name of the resource. To be used in the URL
	Name() string

	// Create a new resource given the json data. Returns
	// the created object or an error
	New(jsonData []byte) ([]byte, error)

	// Get resource for the requested ID
	Fetch(id string) ([]byte, error)

	// Update the resource with the given jsonData.
	// Returns the updated resource or an error
	Update(id string, jsonData []byte) ([]byte, error)

	// Deletes the resource from the DB
	Remove(id string) error
}

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
		result, err := resource.Fetch(vars["id"])
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

		result, err := resource.New(body)
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
		err := resource.Remove(vars["id"])
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
