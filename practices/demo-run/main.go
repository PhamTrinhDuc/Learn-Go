package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{message: "Hello from Go!"}`)
	})
	fmt.Println("Server chạy tại http://localhost:8087")
	http.ListenAndServe(":8087", nil)
}
