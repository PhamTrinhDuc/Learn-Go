package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type User struct {
	ID       int               `json:"id"`
	Username string            `json:"username"`
	Roles    []string          `json:"roles"`
	Metadata map[string]string `json:"metadata"`
}

func GetMockUser() User {
	return User{
		ID:       1,
		Username: "gemini_user",
		Roles:    []string{"admin", "developer"},
		Metadata: map[string]string{
			"location": "Vietnam",
			"status":   "active",
		},
	}
}

func TestJSONFlow() {
	originalUser := GetMockUser()

	// --- 1. Marshal: Chuyển Struct sang JSON ---
	jsonData, err := json.MarshalIndent(originalUser, "", "  ")
	if err != nil {
		fmt.Printf("Marshal lỗi: %v\n", err)
		return
	}
	fmt.Println("--- JSON Output ---")
	fmt.Println(string(jsonData))

	// --- 2. Unmarshal: Chuyển JSON ngược lại Struct ---
	var decodedUser User
	err = json.Unmarshal(jsonData, &decodedUser)
	if err != nil {
		fmt.Printf("Unmarshal lỗi: %v\n", err)
		return
	}

	// --- 3. Kiểm tra tính toàn vẹn (Validation) ---
	if reflect.DeepEqual(originalUser, decodedUser) {
		fmt.Println("\n✅ Thành công: Dữ liệu trước và sau khi Marshal giống hệt nhau!")
	} else {
		fmt.Println("\n❌ Thất bại: Dữ liệu bị thay đổi hoặc mất mát.")
	}
}

func main() {
	TestJSONFlow()
}
