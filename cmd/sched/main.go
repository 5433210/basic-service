package main

import (
	"os"
	"strconv"

	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/sched"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	port, _ := strconv.Atoi(os.Args[1])

	svr := sched.Server{
		Port:          uint16(port),
		IpAddr:        "127.0.0.1",
		LogPath:       "./",
		Name:          "sched",
		ConfPath:      "/Users/zhangweili/Desktop/basic-service/configs/sched.yaml",
		StoreEndpoint: []string{"127.0.0.1:2379"},
		StorePoolSize: 10,
	}

	if err := svr.Run(); err != nil {
		log.ErrorLog(errors.NewError(err))
	}
}
