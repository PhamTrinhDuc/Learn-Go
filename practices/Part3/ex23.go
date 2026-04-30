package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	arr := make([]int, 0, n)

	sum := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Scan(&x)
		arr = append(arr, x)
		sum += x
	}

	sub_sum := 0
	flag := false
	for i := 0; i < len(arr); i++ {
		sub_sum += arr[i]
		if sub_sum == sum-sub_sum-2*arr[i] {
			fmt.Println(i + 1)
			flag = true
			break
		}
	}
	if !flag {
		fmt.Println(-1)
	}
}
