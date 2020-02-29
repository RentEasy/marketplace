package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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

func SetupDatabase() Database {
	db, err := gorm.Open("postgres", "postgres://clone1018@localhost/marketplace?sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}

	db.LogMode(true)

	// Migrate the schema
	if err := db.AutoMigrate(&Property{}).Error; err != nil {
		panic(err)
	}

	return Database{
		db: db,
		Properties: &PropertyQueries{
			db: db,
		},
	}
}

type Database struct {
	db         *gorm.DB
	Properties *PropertyQueries
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
