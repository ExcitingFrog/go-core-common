package grpc_gateway

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	grpcProvider "github.com/ExcitingFrog/go-core-common/grpc"
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
	grpc    *grpcProvider.GRpc
}

func NewGataway(config *Config, grpc *grpcProvider.GRpc) *Gataway {
	if config == nil {
		config = NewConfig()
	}

	return &Gataway{
		Config: config,
		grpc:   grpc,
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

	// wait register
	time.Sleep(2 * time.Second)
	if !g.grpc.Running {
		return errors.New("grpc server not running")
	}
	logrus.Info("gateway server listen on ", g.addr)
	if err := g.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	logrus.Info("gateway start success")

	return nil
}

func (g *Gataway) Options() []grpc.DialOption {
	return g.options
}
