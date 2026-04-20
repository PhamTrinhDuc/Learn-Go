package main

import "fmt"

type Query struct {
	Where string
	Limit int
}

func withLimit(n int) func(*Query) {
	return func(q *Query) {
		q.Limit = n
	}
}

func andWhere(cond string) func(*Query) {
	return func(q *Query) {
		if q.Where == "" {
			q.Where = cond
		} else {
			q.Where += " AND " + cond
		}
	}
}

func BuildQuery(
	base Query,
	options ...func(*Query)) Query {

	for _, option := range options {
		option(&base)
	}
	return base
}

func main() {
	query := BuildQuery(Query{}, withLimit(10), andWhere("city=HN"))
	fmt.Println(query)
}
