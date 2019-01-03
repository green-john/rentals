package rentals

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
