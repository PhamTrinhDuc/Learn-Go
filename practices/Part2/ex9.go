
package main

import "fmt"

func main() {
	var a int
	fmt.Scan(&a)
	if a*2 > 100 {
		fmt.Println("High")
	} else {
		fmt.Println("Low")
	}
}
