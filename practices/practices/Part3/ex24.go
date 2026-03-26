package main

import "fmt"

func main() {
	var x int
	fmt.Scan(&x)

	sum := 0
	list := []int{500, 200, 100, 50, 20, 10}
	for value := range list {
		cnt := 0
		for x % value == 0 {
			cnt += 1
			x /= value
		}
		fmt.Print()
	}
}
