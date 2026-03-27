package main

import (
	"fmt"
	"time"
)

func onStart() time.Time {
	return time.Now()
}

func onFinish(start time.Time) time.Duration {
	return time.Now().Sub(start)
}

func fibo_fn(n int) int {
	if n <= 0 {
		return -1
	}
	if n == 1 {
		return 0
	}
	seq := make([]int, n)
	seq[0], seq[1] = 0, 1
	for i := 2; i < n; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	return seq[n-1]
}

func Observe(
	fn func(int) int,
	onStart func() time.Time,
	onFinish func(time.Time) time.Duration,
) func(int) int {

	return func(x int) int {
		start_time := onStart()
		result := fn(x)
		duration := onFinish(start_time)

		fmt.Printf("Thực thi tổng cộng: %s\n", duration)
		return result
	}
}

func main() {
	res := Observe(fibo_fn, onStart, onFinish)(10)
	fmt.Println(res)
}
