package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/handler"
	"github.com/renteasy/marketplace/internal/database"
	"net/http"
)

type GraphQLType struct {
	Type           *graphql.Object
	QueryFields    graphql.Fields
	MutationFields graphql.Fields
}

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	},
)

var db database.Database

func buildSchema() graphql.SchemaConfig {
	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	registers := []GraphQLType{
		propertyRegister,
		rentalRegister,
	}

	for _, reg := range registers {
		for name, config := range reg.QueryFields {
			schemaConfig.Query.AddFieldConfig(name, config)
		}
		for name, config := range reg.MutationFields {
			schemaConfig.Mutation.AddFieldConfig(name, config)
		}
	}

	return schemaConfig
}

func main() {
	db = database.SetupDatabase("postgres://luke@localhost/marketplace?sslmode=disable")

	schema, err := graphql.NewSchema(buildSchema())
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		FormatErrorFn: func(err error) gqlerrors.FormattedError {
			return gqlerrors.FormatError(err)
		},
	})

	http.Handle("/graphql", h)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
