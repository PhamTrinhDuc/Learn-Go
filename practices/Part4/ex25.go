package main

import "fmt"

func addOne(x int) int {
	return x + 1
}

func square(x int) int {
	return x * x
}

func double(x int) int {
	return x * 2
}

func pipeline(funcs ...func(int) int) func(int) int { // ...: Cho phép nhận nhiều tham số cùng kiểu
	solve := func(x int) int {
		res := funcs[0](x)
		for i := 1; i < len(funcs); i++ {
			res = funcs[i](res)
		}
		return res
	}
	return solve
}

type Transformer func(int) int

func pipeline2(funcs ...Transformer) Transformer {
	solve := func(x int) int {
		res := funcs[0](x)
		for i := 1; i < len(funcs); i++ {
			res = funcs[i](res)
		}
		return res
	}
	return solve
}

func main() {
	var x int
	fmt.Scan(&x)

	res := pipeline(addOne, square, double)(x)
	fmt.Println(res)
}
