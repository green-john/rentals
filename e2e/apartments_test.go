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

func TestCreateApartment(t *testing.T) {
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
	realtorId, err := createUser("realtor", "realtor", "realtor", app.Server.Db)
	tst.Ok(t, err)
	_, err = createUser("client", "client", "client", app.Server.Db)
	tst.Ok(t, err)

	t.Run("Create apartment no auth, fail", func(t *testing.T) {
		// Act
		res, err := tst.MakeRequest("POST", serverUrl+"/apartments", "", []byte(""))
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Create apartment with client, fail", func(t *testing.T) {
		token, err := loginWithUser(t, serverUrl, "client", "client")
		tst.Ok(t, err)
		payload := newApartmentPayload("apt1", "desc", realtorId)

		// Act
		res, err := tst.MakeRequest("POST", serverUrl+"/apartments", token, payload)
		tst.Ok(t, err)

		// Assert
		tst.Assert(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))
	})

	t.Run("Create apartment realtor admin, success", func(t *testing.T) {
		for _, user := range []string{"admin", "realtor"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)
			payload := newApartmentPayload("apt1", "desc", realtorId)

			// Act
			res, err := tst.MakeRequest("POST", serverUrl+"/apartments", token, payload)
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
			rawContent, err := ioutil.ReadAll(res.Body)
			tst.Ok(t, err)

			var apartmentResponse rentals.Apartment

			err = json.Unmarshal(rawContent, &apartmentResponse)

			tst.Assert(t, apartmentResponse.ID >= 1, "Expected id greater than 0")
			tst.Assert(t, apartmentResponse.Name == "apt1", "Got name different name")
			tst.Assert(t, apartmentResponse.RealtorId == realtorId, "Got unexpected realtor")
			tst.Assert(t, apartmentResponse.Available, "Expected apartment to be available")
		}
	})
}

func TestReadApartments(t *testing.T) {
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
	realtorId, err := createUser("realtor", "realtor", "realtor", app.Server.Db)
	tst.Ok(t, err)
	_, err = createUser("client", "client", "client", app.Server.Db)
	tst.Ok(t, err)

	create10Apartments(t, realtorId, app.Server.Db)

	t.Run("Read apartments non auth user", func(t *testing.T) {
		for _, url := range []string{"/apartments", "/apartments/2"} {
			// Act
			res, err := tst.MakeRequest("GET", serverUrl+url, "", []byte(""))
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusUnauthorized,
				fmt.Sprintf("Expected 401, got %d", res.StatusCode))
		}
	})

	t.Run("Read one apartment client, realtor, admin, success", func(t *testing.T) {
		for _, user := range []string{"client", "realtor", "admin"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Act
			res, err := tst.MakeRequest("GET", serverUrl+"/apartments/1", token, []byte(""))
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var returnedApartments rentals.Apartment
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&returnedApartments)
			tst.Ok(t, err)

			tst.Assert(t, returnedApartments.ID == 1,
				fmt.Sprintf("Expected id 1, got %d", returnedApartments.ID))
		}
	})

	t.Run("Read all apartments client, realtor, admin, success", func(t *testing.T) {
		for _, user := range []string{"client", "realtor", "admin"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Act
			res, err := tst.MakeRequest("GET", serverUrl+"/apartments", token, []byte(""))
			tst.Ok(t, err)

			// Assert
			tst.Assert(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var returnedApartments []*rentals.Apartment
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&returnedApartments)
			tst.Ok(t, err)

			tst.Assert(t, len(returnedApartments) == 10,
				fmt.Sprintf("Expected 13 users, got %d", len(returnedApartments)))
		}
	})
}

func newApartmentPayload(name, desc string, realtorId uint) []byte {
	return []byte(fmt.Sprintf(
		`{
"name":"%s",
"description": "%s",
"floorAreaMeters": 50.0,
"pricePerMonthUSD": 500.0,
"roomCount": 4,
"latitude": 41.761536,
"longitude": 12.315237,
"available": true,
"realtorId": %d}`, name, desc, realtorId))
}

func create10Apartments(t *testing.T, realtorId uint, db *gorm.DB) {
	for i := 0; i < 10; i++ {
		user := fmt.Sprintf("user%d", i)
		_, err := createApartment(user, user, realtorId, db)
		tst.Ok(t, err)
	}
}

// Creates a user. Returns its id.
func createApartment(name, desc string, realtorId uint, db *gorm.DB) (uint, error) {
	apartmentResource := &rentals.ApartmentResource{Db: db}

	apartmentData := newApartmentPayload(name, desc, realtorId)
	jsonData, err := apartmentResource.Create([]byte(apartmentData))
	if err != nil {
		return 0, err
	}

	var apartmentId struct {
		Id uint `json:"id"`
	}

	err = json.Unmarshal(jsonData, &apartmentId)
	if err != nil {
		return 0, err
	}

	return apartmentId.Id, nil
}
