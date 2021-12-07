package main

import (
	authz "wailik.com/internal/authz"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

var (
	regoPath = "/Users/zhangweili/Desktop/rbac/data/rbac.rego"
	dataPath = "/Users/zhangweili/Desktop/rbac/data/data.json"
	dbPath   = "/Users/zhangweili/Desktop/rbac/data/db"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	svr := authz.Server{
		Port:     9001,
		IpAddr:   "0.0.0.0",
		RegoPath: regoPath,
		DataPath: dataPath,
		DBPath:   dbPath,
		LogPath:  "",
	}

	if err := svr.Run(); err != nil {
		log.ErrorLog(errors.NewError(err))
	}
}
