package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/lwzm/gqpg"
)

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    gqpg.QueryType,
		Mutation: gqpg.MutationType,
	})

	if err != nil {
		log.Fatal(err)
	}

	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true,
	})

	http.Handle("/", graphqlHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
