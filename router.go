package rentals

import (
	"github.com/gorilla/mux"
	"net/http"
)

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

