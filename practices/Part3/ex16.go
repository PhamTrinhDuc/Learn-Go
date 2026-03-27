package main

import "fmt"

func main() {
	var n int
	arr := make([]int, 0, n)

	fmt.Scan(&n)
	sum_truth := 0
	for i := 0; i < n-1; i++ {
		var x int
		fmt.Scan(&x)
		sum_truth += x
		arr = append(arr, x)
	}

	sum := n * (n + 1) / 2
	fmt.Println(sum - sum_truth)

}
