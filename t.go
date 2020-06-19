package gqpg

import (
	"strings"

	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

var fooType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Foo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var helloField = &graphql.Field{
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		"n": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// fmt.Println(p.Args)
		n, _ := p.Args["n"].(int)
		return "world" + strings.Repeat("!", n), nil
	},
}

var testField = &graphql.Field{
	Type: fooType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return map[string]int{"id": 1}, nil
	},
}

// QueryType exported
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"hello": helloField,
		"test":  testField,
		"fsm":   fsmField,
	},
})

// MutationType exported
var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"hello": helloField,
	},
})
