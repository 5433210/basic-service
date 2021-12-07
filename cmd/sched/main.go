package main

import (
	"wailik.com/internal/captcha"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	svr := captcha.Server{
		Port:    3000,
		IpAddr:  "127.0.0.1",
		LogPath: "./",
	}

	if err := svr.Run(); err != nil {
		log.ErrorLog(errors.NewError(err))
	}
}
