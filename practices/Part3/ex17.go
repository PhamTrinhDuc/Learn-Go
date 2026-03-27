package main

import "fmt"

func main() {
	var n int

	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Scan(&x)
		for j := 0; j < x; j++ {
			fmt.Printf("*")
		}
		fmt.Println()
	}
}
