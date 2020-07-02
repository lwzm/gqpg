package gqpg

import (
	"database/sql"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	pg, err := sql.Open("postgres", "postgres://postgres:x@tyio.net/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db = pg
	// row := db.QueryRow("select id, state, ts, data from fsm where id < $1 limit 1", 2)
}

var fsmType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Fsm",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"state": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"ts": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

var fsmArgs = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"idLte": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"idGte": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"state": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"stateLte": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"stateGte": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"tsLt": &graphql.ArgumentConfig{
		Type: graphql.DateTime,
	},
	"tsGt": &graphql.ArgumentConfig{
		Type: graphql.DateTime,
	},
}

var fsmField = &graphql.Field{
	Type: fsmType,
	Args: fsmArgs,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		clause, args := buildQueryOne(p.Args)
		row := db.QueryRow("select id, state, ts from fsm"+clause, args...)
		var i int
		var s string
		var t time.Time
		if err := row.Scan(&i, &s, &t); err != nil {
			return nil, err
		}
		return object{
			"id":    i,
			"state": s,
			"ts":    t,
		}, nil
	},
}

var fsmsField = &graphql.Field{
	Type: graphql.NewList(fsmType),
	Args: withPage(fsmArgs),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		clause, args := buildQueryWithPage(p.Args)
		rows, err := db.Query("select id, state, ts from fsm"+clause, args...)
		if err != nil {
			return nil, err
		}
		lst := []object{}
		var i int
		var s string
		var t time.Time
		for rows.Next() {
			if err := rows.Scan(&i, &s, &t); err != nil {
				return nil, err
			}
			lst = append(lst, object{
				"id":    i,
				"state": s,
				"ts":    t,
			})
		}
		return lst, nil
	},
}

var fsmsCount = &graphql.Field{
	Type: graphql.Int,
	Args: fsmArgs,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		clause, args := buildQuery(p.Args)
		row := db.QueryRow("select count(id) from fsm"+clause, args...)
		i := 0
		if err := row.Scan(&i); err != nil {
			return nil, err
		}
		return i, nil
	},
}
