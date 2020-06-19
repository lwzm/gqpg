package gqpg

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/iancoleman/strcase"
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

type Fsm0 struct {
	ID    int64  `json:"id"`
	State string `json:"name"`
}

type Fsm struct {
	ID    int64
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
var opsMap = map[string]string{
	"gt":   ">",
	"lt":   "<",
	"gte":  ">=",
	"lte":  "<=",
	"like": "like",
}

var fsmField = &graphql.Field{
	Type: fsmType,
	Args: graphql.FieldConfigArgument{
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
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		fmt.Println(p.Args)
		conditions := "true"
		seq := 1
		args := []interface{}{}
		for k, v := range p.Args {
			key := strcase.ToSnake(k)
			op := "="
			if idx := strings.LastIndex(key, "_"); idx > 0 {
				op = opsMap[key[idx+1:]]
				key = key[:idx]
			}
			conditions += fmt.Sprintf(" and %s %s $%d", key, op, seq)
			args = append(args, v)
			seq++
		}
		fmt.Println(conditions)
		row := db.QueryRow("select id, state, ts from fsm where "+conditions+" limit 1", args...)
		fsm := Fsm{}
		if err := row.Scan(&fsm.ID, &fsm.State, &fsm.Ts); err != nil {
			return nil, err
		}

		return fsm, nil
	},
}
