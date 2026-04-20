package main

import (
	"fmt"
	"strings"
)

func lowerMiddleware(x string) string {
	return strings.ToLower(x)
}

func trimMiddleware(x string) string {
	return strings.TrimSpace(x)
}

func authMiddleware(x string) string {
	if strings.Contains(x, "admin") {
		return "403"
	}
	return "200"
}

type AnonyString func(x string) string

func middleware(funcs ...AnonyString) AnonyString {
	solve := func(x string) string {
		if len(funcs) == 0 {
			return x
		}
		res := funcs[0](x)
		for i := 1; i < len(funcs); i++ {
			res = funcs[i](res)
		}
		return res
	}

	return solve
}

func main() {
	var x string
	fmt.Scan(&x)

	res := middleware(authMiddleware)(x)
	fmt.Println(res)

}
