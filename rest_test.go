package rentals

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func (PageResource) Create(jsonData []byte) ([]byte, error) {
	return jsonData, nil
}

func (PageResource) Read(id string) ([]byte, error) {
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (PageResource) Update(id string, jsonData []byte) ([]byte, error) {
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (PageResource) Delete(id string) error {
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

type verifyCall bool

func (c *verifyCall) call(w http.ResponseWriter, r *http.Request) {
	*c = true
}

func TestRouteAdded(t *testing.T) {
	// Arrange
	var called verifyCall = false
	r := NewGorillaRouter()

	// Act
	r.Add("/one", "POST", called.call)
	res := executeRequest(r, "POST", "/one", nil)

	// Assert
	assert(t, res.Code == http.StatusOK, fmt.Sprintf("Expected 200 got %d", res.Code))
	assert(t, bool(called), "Expected called to be true")
}

func executeRequest(router *GorillaRouter, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, body)
	res := httptest.NewRecorder()
	router.Router.ServeHTTP(res, req)

	return res
}
