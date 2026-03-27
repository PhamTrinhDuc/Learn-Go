package main

import "fmt"

func Memoize() func(int) int {
	sum := 0

	return func(x int) int {
		sum += x
		return sum
	}
}

func main() {
	memorize := Memoize()

	memorize(5)
	memorize(5)

	fmt.Println(memorize(10))
}
