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

	svr := &courier.CourierServer{}
	svr, err := server.LoadConfig(constant.ServiceNameCourier, []string{"."}, svr)
	if err != nil {
		log.ErrorLog(errors.NewError(err))
		return
	}

	svr, err = courier.CreateService(svr)
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
