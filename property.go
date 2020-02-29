package marketplace

import "time"

type Property struct {
	Id             int64
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
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
