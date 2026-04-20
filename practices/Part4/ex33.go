package main

import (
	"fmt"
	"time"
)

func test() int {
	// sleep 1s
	// time.Sleep(1 * time.Second)
	return 10000
}

func WithRateLimit(
	fn func() int,
	limit int,
	interval time.Duration) func() {

	timeQueue := make([]time.Time, 0)

	return func() {
		now := time.Now()
		valid := make([]time.Time, 0)
		for _, t := range timeQueue {
			if now.Sub(t) < interval {
				valid = append(valid, t)
			}
		}

		timeQueue = valid
		if len(timeQueue) >= limit {
			fmt.Println("Rate Limit Exceeded")
			return
		}

		timeQueue = append(timeQueue, time.Now())
		fmt.Println("Send request to function")
		fn()
	}
}

func main() {
	limiter := WithRateLimit(test, 2, time.Second)
	// 3 lần liên tiếp trong ~0ms
	limiter() // ✅ chạy
	limiter() // ✅ chạy
	limiter() // ❌ Rate limit exceeded

	// Đợi hết interval rồi gọi lại
	time.Sleep(time.Second)
	limiter() // ✅ chạy lại được
}
