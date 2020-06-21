package gqpg

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/iancoleman/strcase"
)

var opsMap = map[string]string{
	"gt":   ">",
	"lt":   "<",
	"gte":  ">=",
	"lte":  "<=",
	"like": "like",
}

type any = interface{}

func buildQuery(conditions map[string]any) (string, []any) {
	seq := 1
	clause := "true"
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
func buildPage(conditions map[string]any) (s string) {
	for _, k := range []string{"offset", "limit"} {
		if v, ok := conditions[k]; ok {
			delete(conditions, k)
			n := v.(int)
			s += fmt.Sprintf(" %s %d ", k, n)
		}
	}
	return
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
