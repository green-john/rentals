package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
	"reflect"
	"rentals"
	"strconv"
)

var JsonTagsToFilter = map[string]string{
	"floor_area_meters":   getJsonTag(rentals.Apartment{}, "FloorAreaMeters"),
	"price_per_month_usd": getJsonTag(rentals.Apartment{}, "PricePerMonthUsd"),
	"room_count":          getJsonTag(rentals.Apartment{}, "RoomCount"),
}

type dbApartmentService struct {
	Db *gorm.DB
}

func (ar *dbApartmentService) Create(in rentals.ApartmentCreateInput) (*rentals.ApartmentCreateOutput, error) {
	ar.Db.Create(&(in.Apartment))

	return &rentals.ApartmentCreateOutput{Apartment: in.Apartment}, nil
}

func (ar *dbApartmentService) Read(in rentals.ApartmentReadInput) (*rentals.ApartmentReadOutput, error) {
	apartment, err := getApartment(in.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	return &rentals.ApartmentReadOutput{Apartment: *apartment}, nil
}

func (ar *dbApartmentService) Find(input rentals.ApartmentFindInput) (*rentals.ApartmentFindOutput, error) {
	values, err := url.ParseQuery(input.Query)
	if err != nil {
		return nil, err
	}

	tx := ar.Db.New()
	for dbField, jsonTag := range JsonTagsToFilter {
		if v, ok := values[jsonTag]; ok {
			if !ok || len(v) == 0 {
				continue
			}

			// TODO potential for injection here
			tx = tx.Where(fmt.Sprintf("%s = ?", dbField), v[0])
		}
	}

	var apartments []rentals.Apartment
	tx.Find(&apartments)
	return &rentals.ApartmentFindOutput{Apartments: apartments}, nil
}

func (ar *dbApartmentService) Update(input rentals.ApartmentUpdateInput) (*rentals.ApartmentUpdateOutput, error) {
	apartment, err := getApartment(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	if err := updateFields(apartment, input.Data); err != nil {
		return nil, err
	}

	// Save to DB
	if err = ar.Db.Save(&apartment).Error; err != nil {
		return nil, err
	}
	return &rentals.ApartmentUpdateOutput{Apartment: *apartment}, nil
}

func (ar *dbApartmentService) Delete(input rentals.ApartmentDeleteInput) (*rentals.ApartmentDeleteOutput, error) {
	apartment, err := getApartment(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	ar.Db.Delete(&apartment)
	return &rentals.ApartmentDeleteOutput{Message: "success"}, nil
}

func getApartment(id string, db *gorm.DB) (*rentals.Apartment, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var apartment rentals.Apartment
	if err = db.First(&apartment, intId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rentals.NotFoundError
		}
		return nil, err
	}

	return &apartment, nil
}

func updateFields(apartment *rentals.Apartment, data map[string]interface{}) error {
	if v, ok := data["name"]; ok {
		apartment.Name = v.(string)
	}

	if v, ok := data["description"]; ok {
		apartment.Desc = v.(string)
	}

	if v, ok := data["floorAreaMeters"]; ok {
		apartment.FloorAreaMeters = v.(float32)
	}

	if v, ok := data["pricePerMonthUSD"]; ok {
		apartment.PricePerMonthUsd = v.(float32)
	}

	if v, ok := data["roomCount"]; ok {
		apartment.RoomCount = v.(int)
	}

	if v, ok := data["latitude"]; ok {
		apartment.Latitude = v.(float32)
	}

	if v, ok := data["longitude"]; ok {
		apartment.Longitude = v.(float32)
	}

	if v, ok := data["available"]; ok {
		apartment.Available = v.(bool)
	}

	return nil
}

func getJsonTag(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return field.Tag.Get("json")
}

func NewDbApartmentService(db *gorm.DB) *dbApartmentService {
	return &dbApartmentService{Db: db}
}
