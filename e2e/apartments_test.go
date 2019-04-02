package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"rentals"
	"rentals/auth"
	"rentals/postgres"
	"rentals/transport"
	"rentals/tst"
	"sync"
	"testing"
)

type apartmentResponse struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Desc             string  `json:"description"`
	RealtorId        uint    `json:"realtorId"`
	FloorAreaMeters  float32 `json:"floorAreaMeters"`
	PricePerMonthUsd float32 `json:"pricePerMonthUSD"`
	RoomCount        int     `json:"roomCount"`
	Latitude         float32 `json:"latitude"`
	Longitude        float32 `json:"longitude"`
	Available        bool    `json:"available"`
}

func newServer(t *testing.T) (*transport.Server, func()) {
	t.Helper()

	db, err := postgres.ConnectToDB(true)
	tst.Ok(t, err)
	db.AutoMigrate(rentals.DbModels...)

	authN := auth.NewDbAuthnService(db)
	authZ := auth.NewAuthzService()
	apatService := postgres.NewDbApartmentService(db)
	usrService := postgres.NewDbUserService(db)

	srv, err := transport.NewServer(db, authN, authZ, apatService, usrService)
	tst.Ok(t, err)

	return srv, func() {
		db.DropTableIfExists(rentals.DbModels...)
	}
}

func startServer(wg sync.WaitGroup, addr string, srv *transport.Server) {
	go func() {
		defer wg.Done()
		log.Printf("[ERROR] %s", srv.ServeHTTP(addr))
	}()
}

func TestCRUDApartment(t *testing.T) {
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	srv, clean := newServer(t)
	defer clean()

	// Make sure we delete all things after we are done
	serverUrl := fmt.Sprintf("http://%s", addr)

	wg.Add(1)
	startServer(wg, addr, srv)

	_, err := createUser("admin", "admin", "admin", srv.Db)
	tst.Ok(t, err)
	realtorId, err := createUser("realtor", "realtor", "realtor", srv.Db)
	tst.Ok(t, err)
	_, err = createUser("client", "client", "client", srv.Db)
	tst.Ok(t, err)

	t.Run("CRUD apartment no auth, fail", func(t *testing.T) {
		res, err := tst.MakeRequest("POST", serverUrl+"/apartments", "", []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		for _, url := range []string{"/apartments", "/apartments/2"} {
			res, err := tst.MakeRequest("GET", serverUrl+url, "", []byte(""))
			tst.Ok(t, err)

			tst.True(t, res.StatusCode == http.StatusUnauthorized,
				fmt.Sprintf("Expected 401, got %d", res.StatusCode))
		}

		res, err = tst.MakeRequest("PATCH", serverUrl+"/apartments/1", "", []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))

		res, err = tst.MakeRequest("DELETE", serverUrl+"/apartments/1", "", []byte(""))
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusUnauthorized,
			fmt.Sprintf("Expected 401, got %d", res.StatusCode))
	})

	t.Run("Create Update Delete apartment with client, fail", func(t *testing.T) {
		token, err := loginWithUser(t, serverUrl, "client", "client")
		tst.Ok(t, err)
		newApartmentPayload := newApartmentPayload("apt1", "desc", 5, realtorId)

		res, err := tst.MakeRequest("POST", serverUrl+"/apartments", token, newApartmentPayload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))

		res, err = tst.MakeRequest("PATCH", serverUrl+"/apartments/1", token, newApartmentPayload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))

		res, err = tst.MakeRequest("DELETE", serverUrl+"/apartments/1", token, newApartmentPayload)
		tst.Ok(t, err)

		tst.True(t, res.StatusCode == http.StatusForbidden,
			fmt.Sprintf("Expected 403 got %d", res.StatusCode))
	})

	t.Run("CRUD apartment realtor admin, success", func(t *testing.T) {
		for _, user := range []string{"admin", "realtor"} {
			// Get user token
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Create
			payload := newApartmentPayload("apt1", "desc", 5, realtorId)
			res, err := tst.MakeRequest("POST", serverUrl+"/apartments", token, payload)
			tst.Ok(t, err)

			tst.True(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
			rawContent, err := ioutil.ReadAll(res.Body)
			tst.Ok(t, err)

			var aptRes apartmentResponse
			err = json.Unmarshal(rawContent, &aptRes)

			tst.True(t, aptRes.ID >= 1, "Expected id greater than 0")
			tst.True(t, aptRes.Name == "apt1", "Got name different name")
			tst.True(t, aptRes.RealtorId == realtorId, "Got unexpected realtor")
			tst.True(t, aptRes.Available, "Expected apartment to be available")

			// Read
			apartmentUrl := fmt.Sprintf("%s/apartments/%d", serverUrl, aptRes.ID)
			res, err = tst.MakeRequest("GET", apartmentUrl, token, []byte(""))
			tst.Ok(t, err)

			tst.True(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var retApt apartmentResponse
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&retApt)
			tst.Ok(t, err)

			tst.True(t, retApt.ID == aptRes.ID, fmt.Sprintf("Expected id 1, got %d", retApt.ID))

			// Update
			newData := []byte(`{"id": 100, "name": "newName", "description": "newDesc"}`)
			res, err = tst.MakeRequest("PATCH", apartmentUrl, token, newData)
			tst.Ok(t, err)

			tst.True(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var updApt apartmentResponse
			decoder = json.NewDecoder(res.Body)
			err = decoder.Decode(&updApt)
			tst.Ok(t, err)
			tst.True(t, updApt.ID == retApt.ID,
				fmt.Sprintf("Expected id to be %d, got %d", updApt.ID, retApt.ID))
			tst.True(t, updApt.Name == "newName",
				fmt.Sprintf("Expected name to be newName, got %s", updApt.Name))
			tst.True(t, updApt.Desc == "newDesc",
				fmt.Sprintf("Expected name to be newDesc, got %s", updApt.Desc))
			tst.True(t, updApt.FloorAreaMeters == retApt.FloorAreaMeters,
				fmt.Sprintf("Expected floorArea to be %f, got %f",
					retApt.FloorAreaMeters, updApt.FloorAreaMeters))
			tst.True(t, updApt.PricePerMonthUsd == retApt.PricePerMonthUsd,
				fmt.Sprintf("Expected pricePM to be %f, got %f",
					retApt.PricePerMonthUsd, updApt.PricePerMonthUsd))

			// Delete
			res, err = tst.MakeRequest("DELETE", apartmentUrl, token, []byte(""))
			tst.Ok(t, err)
			tst.True(t, res.StatusCode == http.StatusNoContent,
				fmt.Sprintf("Expected 204, got %d", res.StatusCode))
		}
	})
}

func TestReadAllApartmentsAndSearch(t *testing.T) {
	var wg sync.WaitGroup
	const addr = "localhost:8083"
	srv, clean := newServer(t)
	defer clean()

	serverUrl := fmt.Sprintf("http://%s", addr)

	wg.Add(1)
	startServer(wg, addr, srv)

	_, err := createUser("admin", "admin", "admin", srv.Db)
	tst.Ok(t, err)
	realtorId, err := createUser("realtor", "realtor", "realtor", srv.Db)
	tst.Ok(t, err)
	_, err = createUser("client", "client", "client", srv.Db)
	tst.Ok(t, err)

	create10Apartments(t, realtorId, srv.Db)

	t.Run("Read all apartments client, realtor, admin, success", func(t *testing.T) {
		for _, user := range []string{"client", "realtor", "admin"} {
			token, err := loginWithUser(t, serverUrl, user, user)
			tst.Ok(t, err)

			// Act
			res, err := tst.MakeRequest("GET", serverUrl+"/apartments", token, []byte(""))
			tst.Ok(t, err)

			// True
			tst.True(t, res.StatusCode == http.StatusOK,
				fmt.Sprintf("Expected 200, got %d", res.StatusCode))

			var returnedApartments []apartmentResponse
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&returnedApartments)
			tst.Ok(t, err)

			tst.True(t, len(returnedApartments) == 10,
				fmt.Sprintf("Expected 10 apartments, got %d", len(returnedApartments)))
		}
	})

	t.Run("Search apartments by room count", func(t *testing.T) {
		token, err := loginWithUser(t, serverUrl, "client", "client")
		tst.Ok(t, err)

		// Act
		res, err := tst.MakeRequest("GET", serverUrl+"/apartments?roomCount=4", token, []byte(""))
		tst.Ok(t, err)

		// True
		tst.True(t, res.StatusCode == http.StatusOK,
			fmt.Sprintf("Expected 200, got %d", res.StatusCode))

		var returnedApartments []apartmentResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&returnedApartments)
		tst.Ok(t, err)

		tst.True(t, len(returnedApartments) == 5,
			fmt.Sprintf("Expected 5 apartments, got %d", len(returnedApartments)))
	})
}

func newApartmentPayload(name, desc string, roomCount int, realtorId uint) []byte {
	return []byte(fmt.Sprintf(
		`{
"name":"%s",
"description": "%s",
"floorAreaMeters": 50.0,
"pricePerMonthUSD": 500.0,
"roomCount": %d,
"latitude": 41.761536,
"longitude": 12.315237,
"available": true,
"realtorId": %d}`, name, desc, roomCount, realtorId))
}

func create10Apartments(t *testing.T, realtorId uint, db *gorm.DB) {
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("apt%d", i)
		desc := fmt.Sprintf("desc%d", i)
		_, err := createApartment(name, desc, 2, realtorId, db)
		tst.Ok(t, err)
	}

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("apt%d", 5+i)
		desc := fmt.Sprintf("desc%d", 5+i)
		_, err := createApartment(name, desc, 4, realtorId, db)
		tst.Ok(t, err)
	}
}

func createApartment(name, desc string, roomCount int, realtorId uint, db *gorm.DB) (uint, error) {
	apartmentResource := postgres.NewDbApartmentService(db)

	output, err := apartmentResource.Create(
		rentals.ApartmentCreateInput{
			Apartment: rentals.Apartment{
				Name:             name,
				Desc:             desc,
				RoomCount:        roomCount,
				PricePerMonthUsd: 500.0,
				Latitude:         41.761536,
				Longitude:        12.315237,
				RealtorId:        realtorId,
				Available:        true,
			},
		})
	if err != nil {
		return 0, err
	}

	return uint(output.ID), nil
}
