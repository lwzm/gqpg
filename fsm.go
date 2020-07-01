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

type Fsm struct {
	ID    int
	State string
	Ts    time.Time
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

var pageArgs = graphql.FieldConfigArgument{}

var fsmField = &graphql.Field{
	Type: fsmType,
	Args: fsmArgs,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		clause, args := buildQueryOne(p.Args)
		row := db.QueryRow("select id, state, ts from fsm"+clause, args...)
		fsm := Fsm{}
		if err := row.Scan(&fsm.ID, &fsm.State, &fsm.Ts); err != nil {
			return nil, err
		}
		return fsm, nil
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
		lst := []Fsm{}
		fsm := Fsm{}
		for rows.Next() {
			if err := rows.Scan(&fsm.ID, &fsm.State, &fsm.Ts); err != nil {
				return nil, err
			}
			lst = append(lst, fsm)
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
