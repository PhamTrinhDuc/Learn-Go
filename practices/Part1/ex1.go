package main

// Yêu cầu: Nhập vào chiều dài a và chiều rộng b của một hình chữ nhật (kiểu float64).
// Tính và in ra chu vi và diện tích.
// Input: Hai số thực a và b.
// Output: Chu vi và diện tích, cách nhau bởi dấu cách.

import "fmt"

func perimeter(a, b int) int {
	return 2 * (a + b)
}

func square(a, b int) int {
	return a * b
}

func main() {
	var a, b int
	fmt.Scan(&a)
	fmt.Scan(&b)

	pre := perimeter(a, b)
	sq := square(a, b)

	fmt.Println(pre, sq)

}
