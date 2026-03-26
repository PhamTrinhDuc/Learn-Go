package main

import "fmt"

func main() {
	for true {
		var n int
		fmt.Scan(&n)

		if n == 1 {
			fmt.Println("Hello Gopher")
		} else if n == 2 {
			var x int
			fmt.Scan(&x)
			fmt.Println(x * x)
		} else if n == 3 {
			break
		} else {
			fmt.Println("Lệnh không hợp lệ, nhập lại")
		}
	}
}
