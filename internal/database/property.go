package database

import "github.com/jinzhu/gorm"

type Property struct {
	gorm.Model

	Parcel         string
	Address        string
	City           string
	Coordinates    string
	LotSqft        int
	Sqft           int
	State          string
	Zipcode        string
	UseCode        string
	TotalRooms     int
	Basement       string
	Style          string
	Bedrooms       int
	Grade          string
	Stories        int
	FullBaths      int
	HalfBaths      int
	Condition      string
	YearBuilt      int
	Fireplaces     int
	ExteriorFinish string
	HeatingCooling string
	BasementGarage int
	RoofType       string
}

type PropertyQueries struct {
	db *gorm.DB
}

func (query *PropertyQueries) GetPropertyById(id int) (property Property, err error) {
	if err := query.db.First(&property, id).Error; err != nil {
		return property, err
	}

	return property, nil
}

func (query *PropertyQueries) GetProperties() ([]Property, error) {
	var properties []Property
	if err := query.db.Find(&properties).Error; err != nil {
		return properties, err
	}

	return properties, nil
}

func (query *PropertyQueries) CreateProperty(property *Property) error {
	return query.db.Create(&property).Error
}

func (query *PropertyQueries) UpdateProperty(property *Property) error {
	return query.db.Save(&property).Error
}

func (query *PropertyQueries) DeleteProperty(property *Property) error {
	if property.ID == 0 {
		panic("i'm not deleting the entire DB!")
	}

	return query.db.Delete(property).Error
}
