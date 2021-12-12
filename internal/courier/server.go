package courier

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"

	servicev1 "wailik.com/internal/courier/service/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/microservice"
)

type Server struct {
	Name     string
	Port     uint16
	IpAddr   string
	ConfPath string
	LogPath  string
}

func (svr Server) Run() error {
	app := fiber.New()
	srvc, err := servicev1.New()
	if err != nil {
		panic(err)
	}

	route(app, srvc)
	addr := svr.IpAddr + ":" + strconv.Itoa(int(svr.Port))

	node := microservice.ServiceNode{
		Addr:     "http://" + addr,
		Name:     svr.Name,
		UniqueId: constant.DiscoveryPrifex + "/" + constant.ServiceNameCourier + "/" + xid.New().String(),
		PickMode: microservice.SrvcPickModeRandom,
		RunMode:  microservice.SrvcRunModeFair,
	}
	msrvc, err := microservice.New(node, svr.ConfPath, nil, nil)
	if err != nil {
		return err
	}

	srvc.SetMicroService(msrvc)

	msrvc.Run()

	return app.Listen(addr)
}
