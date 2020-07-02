package gqpg

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/iancoleman/strcase"
)

var opsMap = map[string]string{
	"eq":   "=",
	"neq":  "!=",
	"gt":   ">",
	"lt":   "<",
	"gte":  ">=",
	"lte":  "<=",
	"like": "like",
}

type any = interface{}
type object map[string]any

func buildQuery(conditions map[string]any) (string, []any) {
	seq := 1
	clause := " where true"
	args := []any{}
	for k, v := range conditions {
		key := strcase.ToSnake(k)
		op := "="
		if idx := strings.LastIndex(key, "_"); idx > 0 {
			op = opsMap[key[idx+1:]]
			key = key[:idx]
		}
		clause += fmt.Sprintf(" and %s %s $%d", key, op, seq)
		args = append(args, v)
		seq++
	}
	return clause, args
}

func buildQueryOne(conditions map[string]any) (string, []any) {
	clause, args := buildQuery(conditions)
	return clause + " limit 1", args
}

func buildQueryWithPage(conditions map[string]any) (string, []any) {
	page := ""
	for _, k := range []string{"offset", "limit"} {
		if v, ok := conditions[k]; ok {
			delete(conditions, k)
			n := v.(int)
			page += fmt.Sprintf(" %s %d ", k, n)
		}
	}
	clause, args := buildQuery(conditions)
	return clause + page, args

}

func withPage(org graphql.FieldConfigArgument) graphql.FieldConfigArgument {
	a := graphql.FieldConfigArgument{}
	a["offset"] = &graphql.ArgumentConfig{
		Type: graphql.Int,
	}
	a["limit"] = &graphql.ArgumentConfig{
		Type: graphql.Int,
	}
	for k, v := range org {
		a[k] = v
	}
	return a
}
