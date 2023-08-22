package grpc_gateway

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ExcitingFrog/go-core-common/provider"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
)

type Gataway struct {
	provider.IProvider

	addr   string
	Config *Config
	Mux    *runtime.ServeMux
	server *http.Server
}

func NewGataway(config *Config) *Gataway {
	if config == nil {
		config = NewConfig()
	}

	return &Gataway{
		Config: config,
	}
}

func (g *Gataway) Init() error {
	g.Mux = runtime.NewServeMux()
	return nil
}

func (g *Gataway) Run() error {
	g.addr = fmt.Sprintf(":%d", g.Config.GatawayPort)
	g.server = &http.Server{Addr: g.addr, Handler: g.Mux}

	time.Sleep(3 * time.Second)

	logrus.Info("gateway server listen on ", g.addr)
	if err := g.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	logrus.Info("gateway start success")

	return nil
}
