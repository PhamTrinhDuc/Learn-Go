// Yêu cầu: Khai báo 4 biến trong Go mà không gán giá trị đầu kỳ: int, float64, string, bool.
// In ra giá trị mặc định (Zero Value) của chúng trên cùng một dòng.

// Input: Không có.

// Output: Các giá trị mặc định cách nhau bởi dấu cách.

package main

import "fmt"


func main() {
	var i int 
	var f float64
	var s string 
	var b bool

	fmt.Printf("%v %v %q %v\n", i, f, s, b)
}
