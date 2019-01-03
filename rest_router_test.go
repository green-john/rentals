package rentals

import (
	"fmt"
	"net/http"
	"testing"
)

type MockHttpRouter struct {
	routeMethods map[string][]string
}

func NewMockRouter() *MockHttpRouter {
	rec := make(map[string][]string)
	return &MockHttpRouter{routeMethods: rec}
}

func (r *MockHttpRouter) Add(url string, method string, fn func(http.ResponseWriter, *http.Request)) {
	_, ok := r.routeMethods[url]
	if !ok {
		r.routeMethods[url] = make([]string, 0)
	}

	r.routeMethods[url] = append(r.routeMethods[url], method)
}

func (r *MockHttpRouter) Use(middleware func(http.Handler) http.Handler) {}

func (r *MockHttpRouter) checkAdded(url string, method string) bool {
	methods, ok := r.routeMethods[url]
	if !ok {
		return false
	}

	return contains(methods, method)
}

type PageResource struct{}

func (PageResource) Name() string {
	return "pages"
}

func (PageResource) New(jsonData []byte) ([]byte, error) {
	return jsonData, nil
}

func (PageResource) Fetch(id string) ([]byte, error) {
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (PageResource) Update(id string, jsonData []byte) ([]byte, error) {
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (PageResource) Remove(id string) error {
	return nil
}

func TestCreateRoutes(t *testing.T) {
	// Arrange
	router := NewMockRouter()
	myResource := PageResource{}

	// Act
	CreateRoutes(myResource, router)

	// Assert
	assert(t, router.checkAdded("/pages", "POST"), "no POST /pages")
	assert(t, router.checkAdded("/pages/{id:[0-9]+}", "GET"), "no GET /pages")
	assert(t, router.checkAdded("/pages/{id:[0-9]+}", "PATCH"), "no PATCH /pages")
	assert(t, router.checkAdded("/pages/{id:[0-9]+}", "DELETE"), "no DELETE /pages")
}
