package main

import "fmt"

func main() {
	var str string
	fmt.Scan(&str)
	str += "@"

	cnt := 1
	for i := 1; i < len(str); i++ {
		if str[i] == str[i-1] {
			cnt += 1
		} else {
			fmt.Printf("%c%d", str[i-1], cnt)
			cnt = 1
		}
	}

}
