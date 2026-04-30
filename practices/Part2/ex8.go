// Yêu cầu: Nhập vào hai số thực a, b và một ký tự toán tử op (+, -, *, /).
// Thực hiện phép tính. Nếu là phép chia, hãy kiểm tra mẫu số bằng 0.
// Input: a (float64), op (string), b (float64).
// Output: Kết quả phép tính hoặc "Error" nếu chia cho 0 hoặc toán tử không hợp lệ.

// #1	10 + 5	15
// #2	10 / 0	Error
// #3	7 * 32	1

package main

import "fmt"

func caculate_simple(a int, op string, b int) {
	if op == "+" {
		fmt.Println(a + b)
	} else if op == "-" {
		fmt.Println(a - b)
	} else if op == "*" {
		fmt.Println(a * b)
	} else {
		if b == 0 {
			fmt.Println("Error")
		} else {
			fmt.Println(a / b)
		}
	}
}

func main() {
	var a int
	var op string
	var b int

	fmt.Scan(&a)
	fmt.Scan(&op)
	fmt.Scan(&b)

	caculate_simple(a, op, b)
}
