package main

import (
	"fmt"
	"math"
)

func quadratic_equation(a int, b int, c int) {
	if a == 0 {
		if b == 0 {
			fmt.Println("Vô nghiệm")
		} else {
			x1 := c * -1 / b
			fmt.Println(x1)
		}
	} else {
		delta := b*b - 4*a*c
		if delta == 0 {
			x1 := (b * -1) / (2 * a)
			fmt.Println(x1)
		} else if delta < 0 {
			fmt.Println("Vô nghiệm")
		} else {
			x1 := (float64(-b) + math.Sqrt(float64(delta))) / float64(2*a)
			x2 := (float64(-b) - math.Sqrt(float64(delta))) / float64(2*a)
			fmt.Printf("x1 = %d, x2 = %d", x1, x2)
		}
	}
}

func main() {
	var a, b, c int
	fmt.Scan(&a)
	fmt.Scan(&b)
	fmt.Scan(&c)
	quadratic_equation(a, b, c)
}
