package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
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
	s := "1234"
	fmt.Println(s[1:])
	fmt.Println(s[:2])
	http.ListenAndServe(":8080", nil)
}
