package main

import (
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/server"
	"wailik.com/internal/sched"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	svr := &sched.SchedServer{}
	svr, err := server.LoadConfig(constant.ServiceNameSched, []string{"."}, svr)
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}

	svr, err = sched.CreateService(svr)
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}

	err = server.Run(svr)
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}
}
