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

func TestCreateUser(t *testing.T) {
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

	t.Run("Create user with client, realtor, fail", func(t *testing.T) {
		// Create client
		for _, user := range []string{"client", "realtor"} {
			// Create user
			_, err := createUser(user, user, user, app.server.db)
			ok(t, err)

			token, err := loginWithUser(t, serverUrl, user, user)
			ok(t, err)

			// Act
			res, err := makeRequest("POST", serverUrl+"/users", token, payload)
			ok(t, err)

			// Assert
			assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))
		}
	})

	t.Run("Create user with admin, success", func(t *testing.T) {
		// Create and admin
		_, err := createUser("admin", "admin", "admin", app.server.db)
		ok(t, err)
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
			Id       uint   `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}

		err = json.Unmarshal(rawContent, &userResponse)

		assert(t, userResponse.Username == "john",
			fmt.Sprintf("Expected name john, got %s", userResponse.Username))
		assert(t, userResponse.Role == "client",
			fmt.Sprintf("Expected role client, got %s", userResponse.Role))
		assert(t, userResponse.Id != 0, "Id must be different than 0")
	})
}

func TestCreateApartmentByAdmin(t *testing.T) {
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

	_, err = createUser("admin", "admin", "admin", app.server.db)
	ok(t, err)
	realtorId, err := createUser("realtor", "realtor", "realtor", app.server.db)
	ok(t, err)
	_, err = createUser("client", "client", "client", app.server.db)
	ok(t, err)

	t.Run("Create apartment no auth, fail", func(t *testing.T) {
		// Act
		res, err := makeRequest("POST", serverUrl+"/apartments", "", []byte(""))
		ok(t, err)

		// Assert
		assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Create apartment with client, fail", func(t *testing.T) {
		token, err := loginWithUser(t, serverUrl, "client", "client")
		ok(t, err)
		payload := newApartmentPayload(realtorId)

		// Act
		res, err := makeRequest("POST", serverUrl+"/apartments", token, payload)
		ok(t, err)

		// Assert
		assert(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))
	})

	t.Run("Create apartment realtor admin, success", func(t *testing.T) {
		for _, user := range []string{"admin", "realtor"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			ok(t, err)
			payload := newApartmentPayload(realtorId)

			// Act
			res, err := makeRequest("POST", serverUrl+"/apartments", token, payload)
			ok(t, err)

			// Assert
			assert(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
			rawContent, err := ioutil.ReadAll(res.Body)
			ok(t, err)

			var apartmentResponse Apartment

			err = json.Unmarshal(rawContent, &apartmentResponse)

			assert(t, apartmentResponse.ID >= 1, "Expected id greater than 0")
			assert(t, apartmentResponse.Name == "apt1", "Got name different name")
			assert(t, apartmentResponse.RealtorId == realtorId, "Got unexpected realtor")
			assert(t, apartmentResponse.Available, "Expected apartment to be available")
		}
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

// Creates a user. Returns its id.
func createUser(username, pwd, role string, db *gorm.DB) (uint, error) {
	userResource := &UserResource{db}

	userData := fmt.Sprintf(`{"username": "%s", "password": "%s", "role": "%s"}`,
		username, pwd, role)
	jsonData, err := userResource.Create([]byte(userData))
	if err != nil {
		return 0, err
	}

	var userId struct {
		Id uint `json:"id"`
	}

	err = json.Unmarshal(jsonData, &userId)
	if err != nil {
		return 0, err
	}

	return userId.Id, nil
}

func newApartmentPayload(realtorId uint) []byte {
	return []byte(fmt.Sprintf(
		`{
"name":"apt1",
"description": "nice",
"floorAreaMeters": 50.0,
"pricePerMonthUSD": 500.0,
"roomCount": 4,
"latitude": 41.761536,
"longitude": 12.315237,
"available": true,
"realtorId": %d}`, realtorId))
}
