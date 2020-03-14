package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func SetupDatabase(conn string) Database {
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic("failed to connect database")
	}

	db.LogMode(true)

	// Migrate the schema
	if err := db.AutoMigrate(&Property{}).Error; err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&Rental{}).Error; err != nil {
		panic(err)
	}

	return Database{
		db:         db,
		Properties: &PropertyQueries{db: db},
		Rentals:    &RentalQueries{db: db},
	}
}

type Database struct {
	db         *gorm.DB
	Properties *PropertyQueries
	Rentals    *RentalQueries
}
