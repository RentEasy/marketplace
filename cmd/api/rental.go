package main

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/renteasy/marketplace/internal/database"
)

var rentalType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Rental",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					property, ok := p.Source.(database.Property)
					if !ok {
						return nil, errors.New("could not decode Gorm Model")
					}
					return property.Model.ID, nil
				},
			},
		},
	},
)

var rentalRegister = GraphQLType{
	Type: rentalType,
	QueryFields: graphql.Fields{
		"rental": &graphql.Field{
			Type:        propertyType,
			Description: "Get property by id",
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
			Type:        graphql.NewList(propertyType),
			Description: "Get rentals list",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return db.Rentals.GetRentals()
			},
		},
	},
	MutationFields: graphql.Fields{},
}
