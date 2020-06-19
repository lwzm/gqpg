package gqpg

import (
	"fmt"
	"strings"

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
