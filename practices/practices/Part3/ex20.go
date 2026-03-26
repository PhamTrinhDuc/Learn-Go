package main

import (
	"fmt"
	"math"
)

func ptnt(x int) {
	for i := 2; i <= int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			cnt := 0
			for x%i == 0 {
				x /= i
				cnt += 1
			}
			fmt.Printf("%d^%d ", i, cnt)
		}
	}
	if x != 1 {
		fmt.Printf("%d^%d ", x, 1)
	}
}

func main() {
	var x int
	fmt.Scan(&x)
	ptnt(x)
}
