package main

import "fmt"

func main() {

	// 1. Numerical data types
	var i = 10
	fmt.Printf("Type of x: %T\n", i)
	fmt.Printf("i: %d\n", i)

	// 2. String data types
	// one line
	var foo = "Hello world"
	fmt.Printf(foo + "\n")
	// multiple line
	var bio string = `I am statically typed.
										I was designed at Google.`
	fmt.Printf(bio + "\n")

	// 3. boolean
	var is_value = true

	// 4. Operators (!, &&, ==, !=)
	
}
