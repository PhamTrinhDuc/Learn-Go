
package main

import "fmt"

func check_condition_triagle(a int, b int, c int) string {
	if a == b && b == c {
		return "Deu"
	} else {
		if a+b > c && b+c > a && a+c > b {
			return "Thuong"
		}
		return "Invalid"
	}
}

func main() {
	var a, b, c int
	fmt.Scan(&a)
	fmt.Scan(&b)
	fmt.Scan(&c)

	res := check_condition_triagle(a, b, c)
	fmt.Println(res)
}
