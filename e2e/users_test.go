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

func TestCreateUser(t *testing.T) {
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

	t.Run("Create user no auth, fail", func(t *testing.T) {
		// Act
		res, err := tst.MakeRequest("POST", serverUrl+"/users", "", payload)
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Create user with client, realtor, fail", func(t *testing.T) {
		// Create client
		for _, user := range []string{"client", "realtor"} {
			// Create user
			_, err := createUser(user, user, user, app.Server.Db)
			tst.Ok(t, err)

			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Act
			res, err := tst.MakeRequest("POST", serverUrl+"/users", token, payload)
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Expected 403, got %d", res.StatusCode))
		}
	})

	t.Run("Create user with admin, success", func(t *testing.T) {
		// Create and admin
		_, err := createUser("admin", "admin", "admin", app.Server.Db)
		tst.Ok(t, err)
		token, err := loginWithUser(t, serverUrl, "admin", "admin")
		tst.Ok(t, err)

		// Act
		res, err := tst.MakeRequest("POST", serverUrl+"/users", token, payload)
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
		rawContent, err := ioutil.ReadAll(res.Body)
		tst.Ok(t, err)

		var userResponse struct {
			Id       uint   `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}

		err = json.Unmarshal(rawContent, &userResponse)

		tst.Assert(t, userResponse.Username == "john",
			fmt.Sprintf("Expected name john, got %s", userResponse.Username))
		tst.Assert(t, userResponse.Role == "client",
			fmt.Sprintf("Expected role client, got %s", userResponse.Role))
		tst.Assert(t, userResponse.Id != 0, "Id must be different than 0")
	})
}

func TestReadUser(t *testing.T) {
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

	create10Users(t, app.Server.Db)

	t.Run("Read users non auth user", func(t *testing.T) {
		for _, url := range []string{"/users", "/users/2"} {
			// Act
			res, err := tst.MakeRequest("GET", serverUrl+url, "", []byte(""))
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
				fmt.Sprintf("Expected 401, got %d", res.StatusCode))
		}
	})

	t.Run("Read users client, realtor, fail", func(t *testing.T) {
		for _, user := range []string{"client", "realtor"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			for _, url := range []string{"/users", "/users/2"} {
				// Act
				res, err := tst.MakeRequest("GET", serverUrl+url, token, []byte(""))
				tst.Ok(t, err)

				// Assert
				tst.Assert(t, res.StatusCode == http.StatusForbidden,
					fmt.Sprintf("Expected 403, got %d", res.StatusCode))
			}
		}
	})

	t.Run("Read users admin, succeed", func(t *testing.T) {
		token, err := loginWithUser(t, serverUrl, "admin", "admin")
		tst.Ok(t, err)

		// Act
		res, err := tst.MakeRequest("GET", serverUrl+"/users", token, []byte(""))
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusOK,
			fmt.Sprintf("Expected 200, got %d", res.StatusCode))

		var returnedUsers []*rentals.User
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&returnedUsers)
		tst.Ok(t, err)

		tst.Assert(t, len(returnedUsers) == 13, fmt.Sprintf("Expected 13 users, got %d", len(returnedUsers)))
	})

}

func create10Users(t *testing.T, db *gorm.DB) {
	for i := 0; i < 10; i++ {
		user := fmt.Sprintf("user%d", i)
		_, err := createUser(user, user, "client", db)
		tst.Ok(t, err)
	}
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
