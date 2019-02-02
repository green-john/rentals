package transport

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"rentals/tst"
	"testing"
)

type PageResource struct {
	requestsPerMethod map[string]int
}

func NewPageResource() *PageResource {
	reqCount := make(map[string]int)
	reqCount["GET"] = 0
	reqCount["GETALL"] = 0
	reqCount["POST"] = 0
	reqCount["PATCH"] = 0
	reqCount["DELETE"] = 0

	return &PageResource{requestsPerMethod: reqCount}
}

func (p *PageResource) Create(jsonData []byte) ([]byte, error) {
	p.requestsPerMethod["POST"] += 1
	return jsonData, nil
}

func (p *PageResource) Read(id string) ([]byte, error) {
	p.requestsPerMethod["GET"] += 1
	return []byte(fmt.Sprintf("id:%s", id)), nil
}

func (p *PageResource) Find(query string) ([]byte, error) {
	p.requestsPerMethod["GETALL"] += 1
	return []byte("all"), nil
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
	CreateRoutes("pages", myResource, router)

	res := makeMockRequest(t, router, "/pages", "POST", "")
	tst.Assert(t, res.Code == http.StatusCreated, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages", "GET", "")
	tst.Assert(t, res.Code == http.StatusOK, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages/0", "GET", "")
	tst.Assert(t, res.Code == http.StatusOK, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages/0", "PATCH", "")
	tst.Assert(t, res.Code == http.StatusOK, fmt.Sprintf("Response not ok: %d", res.Code))
	res = makeMockRequest(t, router, "/pages/0", "DELETE", "")
	tst.Assert(t, res.Code == http.StatusNoContent, fmt.Sprintf("Response not ok: %d", res.Code))

	// Assert
	reqCount := myResource.requestsPerMethod["POST"]
	tst.Assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["GET"]
	tst.Assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["GETALL"]
	tst.Assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["PATCH"]
	tst.Assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))

	reqCount = myResource.requestsPerMethod["DELETE"]
	tst.Assert(t, reqCount == 1, fmt.Sprintf("Expected 1 post got %d", reqCount))
}

func makeMockRequest(t *testing.T, router *mux.Router, url, method, body string) *httptest.ResponseRecorder {
	byteBody := bytes.NewBuffer([]byte(body))
	request, err := http.NewRequest(method, url, byteBody)
	tst.Ok(t, err)

	recResponse := httptest.NewRecorder()
	router.ServeHTTP(recResponse, request)

	return recResponse
}
