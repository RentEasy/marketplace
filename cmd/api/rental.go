package main

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/renteasy/marketplace/internal/database"
	"time"
)

var rentalType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Rental",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					rental, ok := p.Source.(database.Rental)
					if !ok {
						return nil, errors.New("could not decode Gorm Model")
					}
					return rental.Model.ID, nil
				},
			},
			"property": &graphql.Field{
				Type: propertyType,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					rental, ok := p.Source.(database.Rental)
					if !ok {
						return nil, errors.New("could not decode Gorm Model")
					}
					return rental.Property, nil
				},
			},
			"bedrooms": &graphql.Field{
				Type:        graphql.Int,
				Description: "How many bedrooms the rental has available to the tenant.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if rental, ok := p.Source.(database.Rental); ok {
						return rental.Bedrooms, nil
					}
					return nil, nil
				},
			},
			"bathrooms": &graphql.Field{
				Type:        graphql.Int,
				Description: "How many bathrooms the rental has available to the tenant.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if rental, ok := p.Source.(database.Rental); ok {
						return rental.Bathrooms, nil
					}
					return nil, nil
				},
			},
			"rentDeposit": &graphql.Field{
				Type:        graphql.Float,
				Description: "How much deposit the tenant would be expected to pay.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if rental, ok := p.Source.(database.Rental); ok {
						return rental.RentDeposit, nil
					}
					return nil, nil
				},
			},
			"rentMonthly": &graphql.Field{
				Type:        graphql.Float,
				Description: "How much the rent is in USD monthly.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if rental, ok := p.Source.(database.Rental); ok {
						return rental.RentMonthly, nil
					}
					return nil, nil
				},
			},
			"listingDate": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "When the property went on the market.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if rental, ok := p.Source.(database.Rental); ok {
						return rental.ListingDate, nil
					}
					return nil, nil
				},
			},
		},
	},
)

var rentalRegister = GraphQLType{
	Type: rentalType,
	QueryFields: graphql.Fields{
		"rental": &graphql.Field{
			Type:        rentalType,
			Description: "Get rental by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					return db.Rentals.GetRentalById(id)
				}
				return nil, nil
			},
		},
		"rentals": &graphql.Field{
			Type:        graphql.NewList(rentalType),
			Description: "Get rentals list",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return db.Rentals.GetRentals()
			},
		},
	},
	MutationFields: graphql.Fields{
		"createRental": &graphql.Field{
			Type:        rentalType,
			Description: "Create a new Rental",
			Args: graphql.FieldConfigArgument{
				// Property Values
				"address": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"city": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"state": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"zipcode": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},

				// Rental Values
				"propertyType": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"rentalStatus": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"unit": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"sqft": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"bedrooms": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"bathrooms": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"stories": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"rentDeposit": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"rentMonthly": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Find or create property
				property, err := db.Properties.FirstOrCreate(database.Property{
					Address: params.Args["address"].(string),
					City:    params.Args["city"].(string),
					State:   params.Args["state"].(string),
					Zipcode: params.Args["zipcode"].(string),
				})
				if err != nil {
					return nil, err
				}

				// Create rental
				unit, unitOk := params.Args["unit"].(string)
				sqft, sqftOk := params.Args["sqft"].(int)
				stories, storiesOk := params.Args["stories"].(int)

				rental := database.Rental{
					Property:    property,
					Bedrooms:    params.Args["bedrooms"].(int),
					Bathrooms:   params.Args["bathrooms"].(int),
					RentDeposit: params.Args["rentDeposit"].(float64),
					RentMonthly: params.Args["rentMonthly"].(float64),
					ListingDate: time.Now(),
				}

				if unitOk {
					rental.Unit = unit
				}
				if sqftOk {
					rental.Sqft = sqft
				}
				if storiesOk {
					rental.Stories = stories
				}

				err = db.Rentals.CreateRental(&rental)
				if err != nil {
					return nil, err
				}

				return rental, nil
			},
		},
	},
}
