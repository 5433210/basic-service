package main

import (
	"wailik.com/internal/authn"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	svr := authn.Server{
		Name:     constant.ServiceNameAuthn,
		Port:     3001,
		IpAddr:   "127.0.0.1",
		DSN:      "root:rootroot@tcp(127.0.0.1:3306)/authn?parseTime=true",
		LogPath:  "./",
		ConfPath: "/Users/zhangweili/Desktop/rbac/configs/authn.yaml",
	}

	if err := svr.Run(); err != nil {
		log.ErrorLog(errors.NewError(err))
	}
}
