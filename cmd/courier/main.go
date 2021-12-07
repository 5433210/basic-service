package main

import (
	"wailik.com/internal/courier"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	svr := courier.Server{
		Name:     constant.ServiceNameCourier,
		Port:     3000,
		IpAddr:   "127.0.0.1",
		LogPath:  "./",
		ConfPath: "/Users/zhangweili/Desktop/rbac/configs/courier.yaml",
	}

	if err := svr.Run(); err != nil {
		log.ErrorLog(errors.NewError(err))
	}
}
