package rentals

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PageResource struct {
	requestsPerMethod map[string]int
}

func NewPageResource() *PageResource {
	reqCount := make(map[string]int)
	reqCount["GET"] = 0
	reqCount["POST"] = 0
	reqCount["PATCH"] = 0
	reqCount["DELETE"] = 0

	return &PageResource{requestsPerMethod: reqCount}
}

func (PageResource) Name() string {
	return "pages"
}

func (p *PageResource) Create(jsonData []byte) ([]byte, error) {
	p.requestsPerMethod["POST"] += 1
	return jsonData, nil
}

func (p *PageResource) Read(id string) ([]byte, error) {
	p.requestsPerMethod["GET"] += 1
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (p *PageResource) Update(id string, jsonData []byte) ([]byte, error) {
	p.requestsPerMethod["PATCH"] += 1
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (p *PageResource) Delete(id string) error {
	p.requestsPerMethod["DELETE"] += 1
	return nil
}

func TestCreateRoutes(t *testing.T) {
	// Arrange
	router := mux.NewRouter()
	myResource := NewPageResource()

	// Act
	CreateRoutes(myResource, router)

	res := makeMockRequest(t, router, "/pages", "POST", "")
	assert(t, res.Code == http.StatusCreated, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages/0", "GET", "")
	assert(t, res.Code == http.StatusOK, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages/0", "PATCH", "")
	assert(t, res.Code == http.StatusOK, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages/0", "DELETE", "")
	assert(t, res.Code == http.StatusNoContent, fmt.Sprintf("Response not ok: %d", res.Code))

	// Assert
	reqCount := myResource.requestsPerMethod["POST"]
	assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["GET"]
	assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["PATCH"]
	assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["DELETE"]
	assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))
}

func makeMockRequest(t *testing.T, router *mux.Router, url, method, body string) *httptest.ResponseRecorder {
	byteBody := bytes.NewBuffer([]byte(body))
	request, err := http.NewRequest(method, url, byteBody)
	ok(t, err)

	recResponse := httptest.NewRecorder()
	router.ServeHTTP(recResponse, request)

	return recResponse
}
