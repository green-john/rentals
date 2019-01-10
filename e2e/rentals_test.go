package e2e

import (
	"encoding/json"
	"fmt"
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
		payload := newApartmentPayload(realtorId)

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
			payload := newApartmentPayload(realtorId)

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
