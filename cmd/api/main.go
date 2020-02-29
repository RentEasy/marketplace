package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
	"github.com/renteasy/marketplace"
	"net/http"
)

var propertyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Property",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
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
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
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
						return getPropertyById(db, id)
					}
					return nil, nil
				},
			},
			/* Get (read) product list
			   http://localhost:8080/product?query={list{id,name,info,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(propertyType),
				Description: "Get properties list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return getProperties(db)
				},
			},
		},
	})
//
//var mutationType = graphql.NewObject(graphql.ObjectConfig{
//	Name: "Mutation",
//	Fields: graphql.Fields{
//		/* Create new product item
//		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
//		*/
//		"create": &graphql.Field{
//			Type:        productType,
//			Description: "Create new product",
//			Args: graphql.FieldConfigArgument{
//				"name": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.String),
//				},
//				"info": &graphql.ArgumentConfig{
//					Type: graphql.String,
//				},
//				"price": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.Float),
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//				rand.Seed(time.Now().UnixNano())
//				product := Product{
//					ID:    int64(rand.Intn(100000)), // generate random ID
//					Name:  params.Args["name"].(string),
//					Info:  params.Args["info"].(string),
//					Price: params.Args["price"].(float64),
//				}
//				products = append(products, product)
//				return product, nil
//			},
//		},
//
//		/* Update product by id
//		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
//		*/
//		"update": &graphql.Field{
//			Type:        productType,
//			Description: "Update product by id",
//			Args: graphql.FieldConfigArgument{
//				"id": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.Int),
//				},
//				"name": &graphql.ArgumentConfig{
//					Type: graphql.String,
//				},
//				"info": &graphql.ArgumentConfig{
//					Type: graphql.String,
//				},
//				"price": &graphql.ArgumentConfig{
//					Type: graphql.Float,
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//				id, _ := params.Args["id"].(int)
//				name, nameOk := params.Args["name"].(string)
//				info, infoOk := params.Args["info"].(string)
//				price, priceOk := params.Args["price"].(float64)
//				product := Product{}
//				for i, p := range products {
//					if int64(id) == p.ID {
//						if nameOk {
//							products[i].Name = name
//						}
//						if infoOk {
//							products[i].Info = info
//						}
//						if priceOk {
//							products[i].Price = price
//						}
//						product = products[i]
//						break
//					}
//				}
//				return product, nil
//			},
//		},
//
//		/* Delete product by id
//		   http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
//		*/
//		"delete": &graphql.Field{
//			Type:        productType,
//			Description: "Delete product by id",
//			Args: graphql.FieldConfigArgument{
//				"id": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.Int),
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//				id, _ := params.Args["id"].(int)
//				product := Product{}
//				for i, p := range products {
//					if int64(id) == p.ID {
//						product = products[i]
//						// Remove from product list
//						products = append(products[:i], products[i+1:]...)
//					}
//				}
//
//				return product, nil
//			},
//		},
//	},
//})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		//Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

var db *sql.DB

func main() {
	db = setupDatabase()

	http.HandleFunc("/property", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func setupDatabase() *sql.DB {
	connStr := "postgres://clone1018@localhost/marketplace?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return db
}

func scanToProperty(row *sql.Row, property marketplace.Property) error {
	return row.Scan(&property.Id,
		&property.Parcel,
		&property.Address,
		&property.City,
		&property.Coordinates,
		&property.LotSqft,
		&property.Sqft,
		&property.State,
		&property.Zipcode,
		&property.UseCode,
		&property.TotalRooms,
		&property.Basement,
		&property.Style,
		&property.Bedrooms,
		&property.Grade,
		&property.Stories,
		&property.FullBaths,
		&property.HalfBaths,
		&property.Condition,
		&property.YearBuilt,
		&property.Fireplaces,
		&property.ExteriorFinish,
		&property.HeatingCooling,
		&property.BasementGarage,
		&property.RoofType,
		&property.CreatedAt,
		&property.UpdatedAt)
}

func getPropertyById(db *sql.DB, id int) (property marketplace.Property, err error) {
	row := db.QueryRow("SELECT * from properties where id = $1", id)
	err = scanToProperty(row, property)

	if err != nil {
		return property, nil
	}

	return property, nil
}

func getProperties(db *sql.DB) ([]marketplace.Property, error) {
	var properties []marketplace.Property
	rows, err := db.Query("SELECT * from properties;")
	if err != nil {
		return properties, err
	}

	for rows.Next() {
		var p marketplace.Property
		err = rows.Scan(
			&p.Id,
			&p.Parcel,
			&p.Address,
			&p.City,
			&p.Coordinates,
			&p.LotSqft,
			&p.Sqft,
			&p.State,
			&p.Zipcode,
			&p.UseCode,
			&p.TotalRooms,
			&p.Basement,
			&p.Style,
			&p.Bedrooms,
			&p.Grade,
			&p.Stories,
			&p.FullBaths,
			&p.HalfBaths,
			&p.Condition,
			&p.YearBuilt,
			&p.Fireplaces,
			&p.ExteriorFinish,
			&p.HeatingCooling,
			&p.BasementGarage,
			&p.RoofType,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		properties = append(properties, p)
	}

	if err != nil {
		return properties, nil
	}

	return properties, nil
}
