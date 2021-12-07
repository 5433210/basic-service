package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	"wailik.com/internal/pkg/log"
)

type Server struct {
	IpAddr   string
	Port     uint16
	RegoPath string
	DataPath string
	DBPath   string
	LogPath  string
}

func (svr Server) Run() error {
	app := fiber.New()
	srvc, err := servicev1.NewService(svr.DBPath, svr.RegoPath, svr.DataPath)
	if err != nil {
		panic(err)
	}

	if err = srvc.LoadData(); err != nil {
		panic(err)
	}

	c, err := srvc.Dump("/")
	if err != nil {
		panic(err)
	}
	log.Info("\n" + string(c))

	route(app, srvc)
	addr := svr.IpAddr + ":" + strconv.Itoa(int(svr.Port))

	return app.Listen(addr)
}
