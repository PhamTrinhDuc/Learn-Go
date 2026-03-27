// Yêu cầu: Nhập vào tháng và năm. In ra số ngày của tháng đó.

// Quy tắc năm nhuận: Chia hết cho 400 HOẶC (chia hết cho 4 VÀ không chia hết cho 100).

// Sử dụng switch cho tháng và if/else lồng nhau cho tháng 2.

// Input: Hai số nguyên month, year.

// Output: Số ngày (ví dụ: 28, 29, 30, 31) hoặc Invalid nếu tháng không hợp lệ.

package main

import "fmt"

func caculate_num_day_in_month(year int, month int) {
	switch month {
	case 1:
		fmt.Println(31)
	case 2:
		if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
			fmt.Println(29)
		} else {
			fmt.Println(28)
		}
	case 3:
		fmt.Println(31)
	case 4:
		fmt.Println(30)
	case 5:
		fmt.Println(31)
	case 6:
		fmt.Println(30)
	case 7:
		fmt.Println(31)
	case 8:
		fmt.Println(31)
	case 9:
		fmt.Println(30)
	case 10:
		fmt.Println(31)
	case 11:
		fmt.Println(30)
	case 12:
		fmt.Println(31)
	}
}

func main() {
	var year, month int
	fmt.Scan(&year)
	fmt.Scan(&month)
	caculate_num_day_in_month(month, year)
}
