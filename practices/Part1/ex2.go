package main

// Yêu cầu: Viết chương trình chuyển đổi nhiệt độ từ độ Celsius (C) sang Fahrenheit (F).
// Công thức: $F = C \times \frac{9}{5} + 32$.Lưu ý: Hãy cẩn thận với phép chia số nguyên trong Go.
// Input: Một số thực C.Output:
// Giá trị F (lấy 2 chữ số thập phân).

import "fmt"

func convert_to_fahrenheit(c float64) float64 {
	F := c*9/5 + 32
	return F
}

func main() {
	var c float64
	fmt.Scan(&c)
	F := convert_to_fahrenheit(c)
	fmt.Println(F)

}
