package main

import "fmt"

func main() {
	var a, b = 10, 20
	fmt.Printf("a: %d, b: %d\n", a, b)

	var (
		foo string = "Hello"
		bar string = "World"
	)

	fmt.Printf("foo: %s, bar %s\n", foo, bar)
}

// run: go run 1_var.go
