package utils

import (
	"os"
	"strconv"
)

// GetEnvString lấy giá trị string từ biến môi trường, nếu không có thì dùng mặc định
func GetEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt lấy giá trị int từ biến môi trường
func GetEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetEnvBool lấy giá trị boolean từ biến môi trường (hỗ trợ 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False)
func GetEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
