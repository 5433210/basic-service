package main

import (
	"wailik.com/internal/courier"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/server"
)

func main() {
	defer log.Flush()

	log.Init(log.OptLevel(log.DebugLevel))

	svr, err := courier.New()
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}

	svr, err = server.LoadConfig(constant.ServiceNameCourier, []string{"."}, svr)
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}

	log.Infof("server conf:%+v", svr)

	err = server.Run(svr)
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}
}
