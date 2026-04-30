// Yêu cầu: Đây là một tính năng cực mạnh của Go. Giả sử bạn có một biến i interface{}. Hãy dùng switch i.(type) để kiểm tra kiểu dữ liệu của nó.
// Nếu là int: in Gấp đôi: {i*2}.
// Nếu là string: in Độ dài: {len(i)}.
// Các kiểu khác: in Unknown.
// Lưu ý: Bạn cần ép kiểu (type assertion) để tính toán sau khi xác định được kiểu

package main

import (
	"fmt"
	"unicode"
)

func check_password(password string) string {
	len_pass := len(password)
	cnt_upper := 0
	cnt_lower := 0
	cnt_digit := 0
	cnt_specific := 0
	for _, c := range password {
		if unicode.IsUpper(c) {
			cnt_upper += 1
		} else if unicode.IsLower(c) {
			cnt_lower += 1
		} else if unicode.IsDigit(c) {
			cnt_digit += 1
		} else {
			cnt_specific += 1
		}
	}
	if len_pass >= 8 && cnt_upper > 0 && cnt_lower > 0 && cnt_digit > 0 && cnt_specific > 0 {
		return "Strong"
	} else {
		return "Weak"
	}
}

func main() {
	var input_str string
	fmt.Scan(&input_str)
	res := check_password(input_str)
	fmt.Println(res)

}
