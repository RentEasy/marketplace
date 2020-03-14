package main

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/renteasy/marketplace/internal/database"
)

var propertyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Property",
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
			"address": &graphql.Field{
				Type: graphql.String,
			},
			"city": &graphql.Field{
				Type: graphql.String,
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
			"zipcode": &graphql.Field{
				Type: graphql.String,
			},
			"sqft": &graphql.Field{
				Type: graphql.Int,
			},
			"style": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var propertyRegister = GraphQLType{
	Type: propertyType,
	QueryFields: graphql.Fields{
		"property": &graphql.Field{
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
					return db.Properties.GetPropertyById(id)
				}
				return nil, nil
			},
		},
		"properties": &graphql.Field{
			Type:        graphql.NewList(propertyType),
			Description: "Get properties list",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return db.Properties.GetProperties()
			},
		},
	},
	MutationFields: graphql.Fields{
		"create": &graphql.Field{
			Type:        propertyType,
			Description: "List a new property",
			Args: graphql.FieldConfigArgument{
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
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				property := database.Property{
					Address: params.Args["address"].(string),
					City:    params.Args["city"].(string),
					State:   params.Args["state"].(string),
					Zipcode: params.Args["zipcode"].(string),
				}

				err := db.Properties.CreateProperty(&property)
				if err != nil {
					return nil, err
				}

				return property, nil
			},
		},
		"update": &graphql.Field{
			Type:        propertyType,
			Description: "Update property by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"city": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"state": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"zipcode": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				address, addressOk := params.Args["address"].(string)
				city, cityOk := params.Args["city"].(string)
				state, stateOk := params.Args["state"].(string)
				zipcode, zipcodeOk := params.Args["zipcode"].(string)

				property, err := db.Properties.GetPropertyById(id)
				if err != nil {
					return nil, err
				}

				if addressOk {
					property.Address = address
				}
				if cityOk {
					property.City = city
				}
				if stateOk {
					property.State = state
				}
				if zipcodeOk {
					property.Zipcode = zipcode
				}

				if err := db.Properties.UpdateProperty(&property); err != nil {
					return nil, err
				}

				return property, nil
			},
		},
		"delete": &graphql.Field{
			Type:        propertyType,
			Description: "Delete property by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				property, err := db.Properties.GetPropertyById(id)
				if err != nil {
					return nil, err
				}

				if err := db.Properties.DeleteProperty(&property); err != nil {
					return property, err
				}

				return property, nil
			},
		},
	},
}
