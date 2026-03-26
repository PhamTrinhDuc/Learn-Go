// 🟣 Bài 5: Tính tổng chữ số (Modulo & Division)Yêu cầu: Nhập vào một số nguyên dương có 3 chữ số.
// Tính tổng các chữ số của nó.Ví dụ: 123 -> $1 + 2 + 3 = 6$.
// Input: Một số nguyên n ($100 \le n \le 999$).
// Output: Tổng các chữ số.
package main

import "fmt"

func sum_digits(x int) int {
	sum := 0
	for x > 0 {
		sum += x % 10
		x /= 10
	}
	return sum
}

func main() {
	var x int
	fmt.Scan(&x)
	sum := sum_digits(x)
	fmt.Println(sum)
}
