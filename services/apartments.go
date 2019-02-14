package services

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
	"reflect"
	"rentals"
	"strconv"
)

var NotFoundError = errors.New("entity not found")

type ApartmentCreateInput struct {
	rentals.Apartment
}

type ApartmentCreateOutput struct {
	rentals.Apartment
}

type ApartmentReadInput struct {
	// ID to lookup the apartment
	Id string
}

type ApartmentReadOutput struct {
	rentals.Apartment
}

type ApartmentFindInput struct {
	Query string
}

type ApartmentFindOutput struct {
	Apartments []rentals.Apartment
}

func (o *ApartmentFindOutput) Public() interface{} {
	return o.Apartments
}

type ApartmentUpdateInput struct {
	Id   string
	Data map[string]interface{}
}

type ApartmentUpdateOutput struct {
	rentals.Apartment
}

type ApartmentDeleteInput struct {
	Id string
}

type ApartmentDeleteOutput struct {
	Message string
}

var JsonTagsToFilter = map[string]string{
	"floor_area_meters":   getJsonTag(rentals.Apartment{}, "FloorAreaMeters"),
	"price_per_month_usd": getJsonTag(rentals.Apartment{}, "PricePerMonthUsd"),
	"room_count":          getJsonTag(rentals.Apartment{}, "RoomCount"),
}

type ApartmentService interface {
	Create(ApartmentCreateInput) (*ApartmentCreateOutput, error)
	Read(ApartmentReadInput) (*ApartmentReadOutput, error)
	Find(ApartmentFindInput) (*ApartmentFindOutput, error)
	Update(ApartmentUpdateInput) (*ApartmentUpdateOutput, error)
	Delete(ApartmentDeleteInput) (*ApartmentDeleteOutput, error)
}

type dbApartmentService struct {
	Db *gorm.DB
}

func (ar *dbApartmentService) Create(in ApartmentCreateInput) (*ApartmentCreateOutput, error) {
	ar.Db.Create(&(in.Apartment))

	return &ApartmentCreateOutput{Apartment: in.Apartment}, nil
}

func (ar *dbApartmentService) Read(in ApartmentReadInput) (*ApartmentReadOutput, error) {
	apartment, err := getApartment(in.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	return &ApartmentReadOutput{Apartment: *apartment}, nil
}

func (ar *dbApartmentService) Find(input ApartmentFindInput) (*ApartmentFindOutput, error) {
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
	return &ApartmentFindOutput{Apartments: apartments}, nil
}

func (ar *dbApartmentService) Update(input ApartmentUpdateInput) (*ApartmentUpdateOutput, error) {
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
	return &ApartmentUpdateOutput{Apartment: *apartment}, nil
}

func (ar *dbApartmentService) Delete(input ApartmentDeleteInput) (*ApartmentDeleteOutput, error) {
	apartment, err := getApartment(input.Id, ar.Db)
	if err != nil {
		return nil, err
	}

	ar.Db.Delete(&apartment)
	return &ApartmentDeleteOutput{Message: "success"}, nil
}

func getApartment(id string, db *gorm.DB) (*rentals.Apartment, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var apartment rentals.Apartment
	if err = db.First(&apartment, intId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NotFoundError
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
