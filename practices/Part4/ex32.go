package main

import "fmt"

func devide(a int, b int) (int, error) {
	if b == 0 {
		return -1, fmt.Errorf("Lỗi, không chia được cho 0")
	}
	return a / b, nil
}

func NewCircuitBreaker(
	fn func(int, int) (int, error),
	threshold int) func(int, int) error {

	failedCount := 0
	isOpen := false

	return func(a int, b int) error {
		if isOpen {
			fmt.Println("Circuit Broken")
		}
		res, err := fn(a, b)
		fmt.Println("Kết quả của hàm fn: ", res)

		if err != nil {
			failedCount += 1
			if failedCount >= threshold {
				isOpen = true
			}
			fmt.Println("Failed count: ", failedCount)
		} else {
			failedCount = 0
		}
		return err
	}
}

func main() {
	breaker := NewCircuitBreaker(devide, 3)
	breaker(5, 0)
	breaker(5, 0)
	breaker(5, 0)
	breaker(5, 0)
	breaker(5, 1)
}
