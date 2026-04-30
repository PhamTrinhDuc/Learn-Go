// 🔴 Bài 3: Hoán vị không dùng biến tạm (Bitwise Operator)
// Yêu cầu: Cho hai số nguyên a và b. Hãy hoán đổi giá trị của chúng mà không sử dụng biến thứ ba (biến tạm). Sử dụng toán tử XOR (^).

// Input: Hai số nguyên a, b.

// Output: Giá trị của a và b sau khi hoán đổ

package main

import "fmt"

func permutations(a, b int) (int, int) {
	var tmp = a
	a = b
	b = tmp
	return a, b
}

func main() {
	var a, b int
	fmt.Scan(&a)
	fmt.Scan(&b)

	a, b = permutations(a, b)
	fmt.Println(a, b)
}
