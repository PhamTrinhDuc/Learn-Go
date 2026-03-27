
package main

import "fmt"

func caculate_tax(a int) float64 {
	if a < 11 {
		return 0.0
	} else if a < 20 {
		return 0.40
	} else {
		return 1.90
	}
}

func main() {
	var a int
	fmt.Scan(&a)
	tax := caculate_tax(a)
	fmt.Println(tax)
}
