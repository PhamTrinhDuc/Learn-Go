package main

import (
	"fmt"

	"charm.land/catwalk/pkg/embedded"
)

const (
	crushDefaultCtxWindow = 128000
	crushDefaultMaxTokens = 4096
)

func main() {
	for _, provider := range embedded.GetAll() {
		for _, m := range provider.Models {
			fmt.Println("Name: ", m.Name)
			fmt.Println("ContextWindow: ", m.ContextWindow)
			fmt.Println("DefaultMaxTokens: ", m.DefaultMaxTokens)
		}
	}
}
