package main

import (
	"fmt"
	"utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Cấu hình log theo kiểu Unix (số giây)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().
		Str("user", "PhamTrinhDuc").
		Int("attempt", 1).
		Msg("Người dùng đang thử đăng nhập")

	fmt.Println(utils.SayHi())
}
