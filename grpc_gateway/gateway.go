package grpc_gateway

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ExcitingFrog/go-core-common/provider"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gataway struct {
	provider.IProvider

	addr    string
	Config  *Config
	Mux     *runtime.ServeMux
	server  *http.Server
	options []grpc.DialOption
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
	g.options = append(g.options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return nil
}

func (g *Gataway) Run() error {
	g.addr = fmt.Sprintf(":%d", g.Config.GatawayPort)
	g.server = &http.Server{Addr: g.addr, Handler: g.Mux}

	time.Sleep(3 * time.Second)

	logrus.Info("gateway server listen on %d", g.addr)
	if err := g.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	logrus.Info("gateway start success")

	return nil
}
