package database

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Rental struct {
	gorm.Model

	PropertyID int
	Property   Property

	Unit        string
	Sqft        string
	Bedrooms    int
	Bathrooms   int
	Stories     int
	RentDeposit float32
	RentMonthly float32

	ListingDate time.Time
}

type RentalQueries struct {
	db *gorm.DB
}

func (query *RentalQueries) GetRentalById(id int) (property Rental, err error) {
	if err := query.db.First(&property, id).Error; err != nil {
		return property, err
	}

	return property, nil
}

func (query *RentalQueries) GetRentals() ([]Rental, error) {
	var properties []Rental
	if err := query.db.Find(&properties).Error; err != nil {
		return properties, err
	}

	return properties, nil
}

func (query *RentalQueries) CreateRental(property *Rental) error {
	return query.db.Create(&property).Error
}

func (query *RentalQueries) UpdateRental(property *Rental) error {
	return query.db.Save(&property).Error
}

func (query *RentalQueries) DeleteRental(property *Rental) error {
	if property.ID == 0 {
		panic("i'm not deleting the entire DB!")
	}

	return query.db.Delete(property).Error
}
