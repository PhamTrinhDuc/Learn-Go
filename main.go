package main

import (
	"fmt"
	"time"
)

func main() {
	value := int(time.Duration.Seconds())
	fmt.Println(value)
}
