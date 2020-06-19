package gqpg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildQuery(t *testing.T) {
	a := assert.New(t)
	conditions := map[string]any{
		"id":       8,
		"ageGt":    10,
		"nameLike": "foo",
	}
	clause, args := buildQuery(conditions)
	a.Contains(clause, "and id = $")
	a.Contains(clause, "and age > $")
	a.Contains(clause, "and name like $")
	a.Contains(args, "foo")
}
