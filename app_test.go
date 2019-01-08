package rentals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
)

func TestCreateUserByAdmin(t *testing.T) {
	// Arrange
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	app, err := NewApp(addr)
	ok(t, err)
	ok(t, app.Setup())

	// Make sure we delete all things after we are done
	defer app.dropDB()
	serverUrl := fmt.Sprintf("http://%s", addr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[ERROR] %s", app.ServeHTTP())
	}()

	payload := []byte(`{"username":"john", "password": "secret", "role": "client"}`)

	t.Run("Create user no auth, fail", func(t *testing.T) {
		// Act
		res, err := makeRequest("POST", serverUrl+"/users", "", payload)
		ok(t, err)

		// Assert
		assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Create user with admin, success", func(t *testing.T) {
		// Create and admin
		ok(t, createAdmin("admin", "admin", app.server.db))
		token, err := loginWithUser(t, serverUrl, "admin", "admin")
		ok(t, err)

		// Act
		res, err := makeRequest("POST", serverUrl+"/users", token, payload)
		ok(t, err)

		// Assert
		assert(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
		rawContent, err := ioutil.ReadAll(res.Body)
		ok(t, err)

		var userResponse struct {
			Username string `json:"username"`
			Role     string `json:"role"`
		}

		err = json.Unmarshal(rawContent, &userResponse)

		assert(t, userResponse.Username == "john",
			fmt.Sprintf("Expected name john, got %s", userResponse.Username))
		assert(t, userResponse.Role == "client",
			fmt.Sprintf("Expected role client, got %s", userResponse.Role))
	})
}

func loginWithUser(t *testing.T, serverUrl, username, pwd string) (string, error) {
	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, pwd)
	response, err := makeRequest("POST", serverUrl+"/login", "", []byte(body))
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 399 {
		t.Errorf("Got %d code, expected 2XX", response.StatusCode)
	}

	var token struct {
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&token)

	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func makeRequest(method, url, authToken string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set json as default
	req.Header.Set("Content-Type", "application/json")

	// If an auth token is provided, add it to the headers
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	return http.DefaultClient.Do(req)
}

// Hacky way of quickly creating an admin
func createAdmin(username, pwd string, db *gorm.DB) error {
	userResource := &UserResource{db}

	userData := fmt.Sprintf(`{"username": "%s", "password": "%s", "role": "admin"}`,
		username, pwd)
	_, err := userResource.Create([]byte(userData))
	if err != nil {
		return err
	}

	return nil
}
