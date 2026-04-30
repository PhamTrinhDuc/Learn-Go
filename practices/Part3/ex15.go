package main

import "fmt"

func main() {
	var p float64
	var r float64
	var n int

	fmt.Scan(&p)
	fmt.Scan(&r)
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		p = p + p*(r/100)
		fmt.Printf("Năm %d: %.2f\n", i+1, p)
	}
}
