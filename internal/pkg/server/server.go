package server

import (
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
)

type IServer interface {
	GetName() string
	GetPort() int
	GetIpAddr() string
	Bind(*fiber.App)
	GetMicroServiceConfig() microservice.MicroServiceConfig
	SetMicroService(microservice.MicroService)
}

type Server struct {
	Name               string
	Port               uint16
	IpAddr             string
	MicroServiceConfig microservice.MicroServiceConfig
}

func (s *Server) GetName() string {
	return s.Name
}
func (s *Server) GetPort() int {
	return int(s.Port)
}

func (s *Server) GetIpAddr() string {
	return s.IpAddr
}

func (s *Server) GetMicroServiceConfig() microservice.MicroServiceConfig {
	return s.MicroServiceConfig
}

func (s *Server) SetMicroService(ms microservice.MicroService) {

}

func (s *Server) Bind(*fiber.App) {

}

func LoadConfig[T any](serviceName string, configPaths []string, config *T) (*T, error) {
	log.Info("new Server...")
	viper.SetConfigName(serviceName)
	for _, path := range configPaths {
		log.Infof("name:%+v, path:%+v", serviceName, path)
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	log.Info(viper.GetString("name"))
	log.Info(viper.GetString("port"))
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	log.Infof("config:%+v", config)
	return config, nil
}

func Run[T IServer](server T) error {
	app := fiber.New()
	server.Bind(app)
	addr := server.GetIpAddr() + ":" + strconv.Itoa(server.GetPort())

	node := microservice.ServiceNode{
		Addr:     "http://" + addr,
		Name:     server.GetName(),
		UniqueId: constant.DiscoveryPrifex + "/" + server.GetName() + "/" + xid.New().String(),
		PickMode: microservice.SrvcPickModeRandom,
		RunMode:  microservice.SrvcRunModeFair,
	}

	msrvc, err := microservice.New(node, server.GetMicroServiceConfig(), nil, nil)
	if err != nil {
		return err
	}

	server.SetMicroService(msrvc)
	msrvc.Start()
	return app.Listen(addr)
}
