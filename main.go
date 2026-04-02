package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "123")
	ctx = context.WithValue(ctx, "tenant_id", "234")
	fmt.Println(ctx.Value("user_id"))

}
