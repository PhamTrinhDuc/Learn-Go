package main

import "fmt"

func main() {
	pi := 0.0
	var n int
	fmt.Scan(&n)

	mau := -1.0
	for i := 0; i < n; i++ {
		mau += 2.0
		if i%2 == 0 {
			pi += 1.0 / mau
		} else {
			pi -= 1.0 / mau
		}
	}

	pi *= 4.0
	fmt.Printf("%.5f", pi)

}
