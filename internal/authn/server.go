package authn

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"

	servicev1 "wailik.com/internal/authn/service/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
)

type Server struct {
	Name     string
	Port     uint16
	IpAddr   string
	DSN      string
	ConfPath string
	LogPath  string
}

func (svr Server) Run() error {
	app := fiber.New()
	srvc, err := servicev1.New(svr.DSN)
	if err != nil {
		panic(err)
	}

	route(app, srvc)
	addr := svr.IpAddr + ":" + strconv.Itoa(int(svr.Port))

	node := microservice.ServiceNode{
		Addr:     "http://" + addr,
		Name:     svr.Name,
		UniqueId: constant.DiscoveryPrifex + "/" + constant.ServiceNameAuthn + "/" + xid.New().String(),
		RunMode:  microservice.SrvcRunModeFair,
		PickMode: microservice.SrvcPickModeHash,
	}
	msrvc, err := microservice.New(node, svr.ConfPath, nil, nil)
	if err != nil {
		log.Debugf("new micro service error")

		return err
	}

	srvc.SetMicroService(msrvc)

	msrvc.Run()

	return app.Listen(addr)
}
