package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"rentals"
	"rentals/tst"
	"testing"
)

func TestFindApartment(t *testing.T) {
	// Arrange
	db, err := ConnectToDB(true)
	tst.Ok(t, err)

	db.AutoMigrate(rentals.DbModels...)
	defer db.DropTableIfExists(rentals.DbModels...)

	aptResource := &dbApartmentService{Db: db}

	createRealtor(t, db)
	createApartments(t, aptResource)

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
			res, err := aptResource.Find(rentals.ApartmentFindInput{Query: elt.query})
			tst.Ok(t, err)

			//err = json.Unmarshal(res, &retApts)
			//tst.Ok(t, err)

			for idx, apt := range res.Apartments {
				tst.True(t, apt.Name == elt.resultIds[idx],
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
func createApartments(t *testing.T, s *dbApartmentService) {
	for area := 1; area <= 2; area++ {
		for price := 1; price <= 2; price++ {
			for rooms := 1; rooms <= 2; rooms++ {
				name := fmt.Sprintf("%d|%d|%d", area, price, rooms)
				payload := newApartmentPayload(name, name, float32(area), float32(price), rooms, 1)
				_, err := s.Create(payload)
				tst.Ok(t, err)
			}
		}
	}
}

func newApartmentPayload(name, desc string, area, price float32, roomCount int,
	realtorId uint) rentals.ApartmentCreateInput {
	return rentals.ApartmentCreateInput{
		Apartment: rentals.Apartment{
			Name:             name,
			Desc:             desc,
			FloorAreaMeters:  area,
			PricePerMonthUsd: price,
			RoomCount:        roomCount,
			RealtorId:        realtorId,
			Longitude:        34.3222223,
			Latitude:         21.233449,
			Available:        true,
		},
	}
}

func createRealtor(t *testing.T, db *gorm.DB) {
	usrService := NewDbUserService(db)

	_, err := usrService.Create(rentals.UserCreateInput{
		Username: "user",
		Password: "pass",
		Role:     "realtor",
	})
	tst.Ok(t, err)
}
