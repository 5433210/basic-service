package captcha

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/captcha/service/v1"
)

type Server struct {
	Port    uint16
	IpAddr  string
	LogPath string
}

func (svr Server) Run() error {
	app := fiber.New()
	srvc, err := servicev1.New()
	if err != nil {
		panic(err)
	}

	route(app, srvc)
	addr := svr.IpAddr + ":" + strconv.Itoa(int(svr.Port))

	return app.Listen(addr)
}
