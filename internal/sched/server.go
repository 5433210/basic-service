package sched

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
	servicev1 "wailik.com/internal/sched/service/v1"
)

type Server struct {
	Port          uint16
	IpAddr        string
	LogPath       string
	Name          string
	ConfPath      string
	StoreEndpoint []string
	StorePoolSize int
}

func (svr Server) Run() error {
	log.Debug("server running...")
	app := fiber.New(fiber.Config{})
	srvc, err := servicev1.New(svr.StoreEndpoint, svr.StorePoolSize)
	if err != nil {
		panic(err)
	}

	route(app, srvc)
	addr := svr.IpAddr + ":" + strconv.Itoa(int(svr.Port))

	node := microservice.ServiceNode{
		Addr:     "http://" + addr,
		Name:     svr.Name,
		UniqueId: constant.DiscoveryPrifex + "/" + constant.ServiceNameSched + "/" + xid.New().String(),
		PickMode: microservice.SrvcPickModeMaster,
		RunMode:  microservice.SrvcRunModeMasterSlave,
	}

	log.Debugf("%v", node)
	msrvc, err := microservice.New(node, svr.ConfPath)
	if err != nil {
		return err
	}

	srvc.SetMicroService(msrvc)

	msrvc.Run()

	return app.Listen(addr)
}
