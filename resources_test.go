package rentals

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"tournaments/tst"
)

func TestFindApartment(t *testing.T) {
	// Arrange
	db, err := ConnectToDB()
	tst.Ok(t, err)

	db.AutoMigrate(DbModels...)
	defer db.DropTableIfExists(DbModels...)

	aptResource := &ApartmentResource{Db: db}

	tst.Ok(t, createRealtor(db))
	tst.Ok(t, createApartments(aptResource))

	for _, elt := range []struct {
		query     string
		resultIds []string
	}{
		{"", []string{"1|1|1", "1|1|2", "1|2|1", "1|2|2", "2|1|1", "2|1|2", "2|2|1", "2|2|2"}},
		{"floorAreaMeters=1", []string{"1|1|1", "1|1|2", "1|2|1", "1|2|2"}},
		{"floorAreaMeters=1&pricePerMonthUSD=1", []string{"1|1|1", "1|1|2"}},
		{"floorAreaMeters=1&pricePerMonthUSD=1&roomCount=1", []string{"1|1|1"}},
		{"floorAreaMeters=1&pricePerMonthUSD=2", []string{"1|2|1", "1|2|2"}},
		{"pricePerMonthUSD=2&roomCount=1", []string{"1|2|1", "2|2|1"}},
	} {
		t.Run(fmt.Sprintf("%s -> %v", elt.query, elt.resultIds), func(t *testing.T) {
			// Act
			var retApts []Apartment
			res, err := aptResource.Find(elt.query)
			tst.Ok(t, err)

			err = json.Unmarshal(res, &retApts)
			tst.Ok(t, err)

			for idx, apt := range retApts {
				tst.Assert(t, apt.Name == elt.resultIds[idx],
					fmt.Sprintf("Expected %s, got %s", elt.resultIds[idx], apt.Name))
			}
		})
	}
}

// Creates 8 apartments with the following attributes
//  area, price, roomCount
//    1      1      1
//    1      1      2
//    1      2      1
//    1      2      2
//    2      1      1
//    2      1      2
//    2      2      1
//    2      2      2
func createApartments(r *ApartmentResource) error {
	for area := 1; area <= 2; area++ {
		for price := 1; price <= 2; price++ {
			for rooms := 1; rooms <= 2; rooms++ {
				name := fmt.Sprintf("%d|%d|%d", area, price, rooms)
				payload := newApartmentPayload(name, name, float32(area), float32(price), rooms, 1)
				_, err := r.Create(payload)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func newApartmentPayload(name, desc string, area, price float32, roomCount int, realtorId uint) []byte {
	return []byte(fmt.Sprintf(
		`{
"name":"%s",
"description": "%s",
"floorAreaMeters": %f,
"pricePerMonthUSD": %f,
"roomCount": %d,
"latitude": 41.761536,
"longitude": 12.315237,
"available": true,
"realtorId": %d}`, name, desc, area, price, roomCount, realtorId))
}

func createRealtor(db *gorm.DB) error {
	usrResource := &UserResource{Db: db}

	usrData := []byte(`{"username":"user", "password": "pass", "role": "realtor"}`)
	_, err := usrResource.Create(usrData)

	if err != nil {
		return err
	}

	return nil
}
