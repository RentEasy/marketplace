package main

import (
	"database/sql"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/lib/pq"
	"os"
)

type TaxParcelInput struct {
	ObjectId12 int `csv:"OBJECTID_12"`
	ObjectId int `csv:"OBJECTID"`
	RevisionDate string `csv:"REVISIONDATE"`
	CvtTaxCode string `csv:"CVTTAXCODE"`
	CvtTaxDescription string `csv:"CVTTAXDESCRIPTION"`
	Pin string `csv:"PIN"`
	SiteAddress string `csv:"SITEADDRESS"`
	SiteCity string `csv:"SITECITY"`
	SiteState string `csv:"SITESTATE"`
	SiteZip5 int `csv:"SITEZIP5"`
	AssessedValue float32 `csv:"ASSESSEDVALUE"`
	TaxableValue float32 `csv:"TAXABLEVALUE"`
	NumBeds int `csv:"NUM_BEDS"`
	NumBaths int `csv:"NUM_BATHS"`
	StructureDesc string `csv:"STRUCTURE_DESC"`
	LivingAreaSqft int `csv:"LIVING_AREA_SQFT"`
	Shapearea string `csv:"Shapearea"`
	Shapelen string `csv:"Shapelen"`
}

func main() {
	connStr := "postgres://clone1018@localhost/marketplace?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * from properties;")
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)


	file, err := os.Open("/Users/clone1018/Downloads/OC_Tax_Parcels_Public.csv")
	if err != nil {
		panic(err)
	}
	properties := []*TaxParcelInput{}

	if err := gocsv.UnmarshalFile(file, &properties); err != nil { // Load clients from file
		panic(err)
	}

	txn, err := db.Begin()
	if err != nil {
		panic(err)
	}


	stmt, err := txn.Prepare(pq.CopyIn("properties",
		"parcel",
		"address",
		"city",
		"state",
		"zipcode",
		"bedrooms",
		"full_baths",
		"style",
		"sqft",
		))
	if err != nil {
		panic(err)
	}

	for _, property := range properties {
		_, err = stmt.Exec(
			property.Pin,
			property.SiteAddress,
			property.SiteCity,
			property.SiteState,
			property.SiteZip5,
			property.NumBeds,
			property.NumBaths,
			property.StructureDesc,
			property.LivingAreaSqft,
			)
		if err != nil {
			panic(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	err = stmt.Close()
	if err != nil {
		panic(err)
	}

	err = txn.Commit()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d!\n", len(properties))
}