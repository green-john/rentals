package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"tournaments"
	"tournaments/tst"
)

type userResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func TestCRUDUsers(t *testing.T) {
	// Arrange
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	app, err := rentals.NewApp(addr)
	tst.Ok(t, err)
	tst.Ok(t, app.Setup())

	// Make sure we delete all things after we are done
	defer app.DropDB()
	serverUrl := fmt.Sprintf("http://%s", addr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[ERROR] %s", app.ServeHTTP())
	}()

	payload := []byte(`{"username":"john", "password": "secret", "role": "client"}`)

	t.Run("CRUD user no auth, fail", func(t *testing.T) {
		// Create
		res, err := tst.MakeRequest("POST", serverUrl+"/users", "", payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		// Read
		res, err = tst.MakeRequest("GET", serverUrl+"/users", "", payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		res, err = tst.MakeRequest("GET", serverUrl+"/users/1", "", payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		// Update
		res, err = tst.MakeRequest("PATCH", serverUrl+"/users/1", "", payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		// Delete
		res, err = tst.MakeRequest("DELETE", serverUrl+"/users/1", "", payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("CRUD user with client, realtor, fail", func(t *testing.T) {
		for _, user := range []string{"client", "realtor"} {
			// Create and get token
			_, err := createUser(user, user, user, app.Server.Db)
			tst.Ok(t, err)
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Create
			res, err := tst.MakeRequest("POST", serverUrl+"/users", token, payload)
			tst.Ok(t, err)
			tst.Assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))

			// Read
			res, err = tst.MakeRequest("GET", serverUrl+"/users", token, payload)
			tst.Ok(t, err)
			tst.Assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))

			res, err = tst.MakeRequest("GET", serverUrl+"/users/1", token, payload)
			tst.Ok(t, err)
			tst.Assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))

			// Update
			res, err = tst.MakeRequest("PATCH", serverUrl+"/users/1", token, payload)
			tst.Ok(t, err)
			tst.Assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))

			// Delete
			res, err = tst.MakeRequest("DELETE", serverUrl+"/users/1", token, payload)
			tst.Ok(t, err)
			tst.Assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))
		}
	})

	t.Run("CRUD user with admin, success", func(t *testing.T) {
		// Create and admin
		_, err := createUser("admin", "admin", "admin", app.Server.Db)
		tst.Ok(t, err)
		token, err := loginWithUser(t, serverUrl, "admin", "admin")
		tst.Ok(t, err)

		// Create
		res, err := tst.MakeRequest("POST", serverUrl+"/users", token, payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
		rawContent, err := ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var usrRes userResponse
		err = json.Unmarshal(rawContent, &usrRes)
		tst.Assert(t, usrRes.Username == "john",
			fmt.Sprintf("Expected name john, got %s", usrRes.Username))
		tst.Assert(t, usrRes.Role == "client",
			fmt.Sprintf("Expected role client, got %s", usrRes.Role))
		tst.Assert(t, usrRes.ID != 0, "Id must be different than 0")

		// Read
		userUrl := fmt.Sprintf("%s/users/%d", serverUrl, usrRes.ID)
		res, err = tst.MakeRequest("GET", userUrl, token, payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusOK, fmt.Sprintf("Expected 200 got %d", res.StatusCode))
		rawContent, err = ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var retUser userResponse
		err = json.Unmarshal(rawContent, &retUser)
		tst.Assert(t, retUser.Username == usrRes.Username,
			fmt.Sprintf("Expected name %s, got %s", usrRes.Username, retUser.Username))
		tst.Assert(t, retUser.Role == usrRes.Role,
			fmt.Sprintf("Expected role %s, got %s", usrRes.Role, retUser.Role))
		tst.Assert(t, retUser.ID == usrRes.ID,
			fmt.Sprintf("Expected id %d, got %d", usrRes.ID, retUser.ID))

		// Update
		payload = []byte(`{"id":100, "username": "newusr", "role": "realtor"}`)
		res, err = tst.MakeRequest("PATCH", userUrl, token, payload)
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusOK, fmt.Sprintf("Expected 200 got %d", res.StatusCode))
		rawContent, err = ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var updUser userResponse
		err = json.Unmarshal(rawContent, &updUser)
		tst.Assert(t, updUser.Username == "john",
			fmt.Sprintf("Expected name john, got %s", updUser.Username))
		tst.Assert(t, updUser.Role == "realtor",
			fmt.Sprintf("Expected role realtor, got %s", updUser.Role))
		tst.Assert(t, updUser.ID == usrRes.ID,
			fmt.Sprintf("Expected id %d, got %d", usrRes.ID, updUser.ID))

		// Delete
		res, err = tst.MakeRequest("DELETE", userUrl, token, []byte(""))
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusNoContent,
			fmt.Sprintf("Expected 204, got %d", res.StatusCode))

		res, err = tst.MakeRequest("GET", userUrl, token, []byte(""))
		tst.Ok(t, err)
		tst.Assert(t, res.StatusCode == http.StatusNotFound,
			fmt.Sprintf("Expected 404, got %d", res.StatusCode))
	})
}

func TestFetchOwnUserData(t *testing.T) {
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	app, err := rentals.NewApp(addr)
	tst.Ok(t, err)
	tst.Ok(t, app.Setup())

	// Make sure we delete all things after we are done
	defer app.DropDB()
	serverUrl := fmt.Sprintf("http://%s", addr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[ERROR] %s", app.ServeHTTP())
	}()

	_, err = createUser("admin", "admin", "admin", app.Server.Db)
	tst.Ok(t, err)
	_, err = createUser("realtor", "realtor", "realtor", app.Server.Db)
	tst.Ok(t, err)
	_, err = createUser("client", "client", "client", app.Server.Db)
	tst.Ok(t, err)

	t.Run("Can't create user with same username", func(t *testing.T) {
		_, err := createUser("admin", "admin", "admin", app.Server.Db)
		tst.Assert(t, err != nil, "Expected error, got success")
	})

	t.Run("Read profile not logged in", func(t *testing.T) {
		// Act
		res, err := tst.MakeRequest("GET", serverUrl+"/profile", "", []byte(""))
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Read own user data client, realtor, admin", func(t *testing.T) {
		for _, user := range []string{"client", "realtor", "admin"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Act
			res, err := tst.MakeRequest("GET", serverUrl+"/profile", token, []byte(""))
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var returnedUser rentals.User
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&returnedUser)
			tst.Ok(t, err)

			assertUser(t, &returnedUser, user, user)
		}
	})
}

func TestCreateClient(t *testing.T) {
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	app, err := rentals.NewApp(addr)
	tst.Ok(t, err)
	tst.Ok(t, app.Setup())

	// Make sure we delete all things after we are done
	defer app.DropDB()
	serverUrl := fmt.Sprintf("http://%s", addr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[ERROR] %s", app.ServeHTTP())
	}()

	_, err = createUser("admin", "admin", "admin", app.Server.Db)
	tst.Ok(t, err)

	t.Run("Create client, not logged in, error", func(t *testing.T) {
		// Act
		payload := []byte(`{"username": "client1", "password": "client1"}`)
		res, err := tst.MakeRequest("POST", serverUrl+"/newClient", "", payload)
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusCreated,
			fmt.Sprintf("Expected 201, got %d", res.StatusCode))

		var returnedUser rentals.User
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&returnedUser)
		tst.Ok(t, err)

		assertUser(t, &returnedUser, "client1", "client")
	})
}

func assertUser(t *testing.T, user *rentals.User, username, role string, ) {
	tst.Assert(t, user.Username == username,
		fmt.Sprintf("Expected username %s, got %s", username, user.Username))

	tst.Assert(t, user.Role == role,
		fmt.Sprintf("Expected role %s, got %s", role, user.Role))
}

// Creates a user. Returns its id.
func createUser(username, pwd, role string, db *gorm.DB) (uint, error) {
	userResource := &rentals.UserResource{Db: db}

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
