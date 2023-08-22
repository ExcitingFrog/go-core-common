package grpc

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/ExcitingFrog/go-core-common/provider"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRpc struct {
	provider.IProvider

	Server *grpc.Server
	Config *Config
}

func NewGRpc(config *Config) *GRpc {
	if config == nil {
		config = NewConfig()
	}
	return &GRpc{
		Config: config,
	}
}

func recoverHandler(_ context.Context, p interface{}) error {
	err := status.Errorf(codes.Internal, "%v", p)
	logrus.WithError(err).Errorf("[gRPC] Service panic, stack: \n%s", debug.Stack())
	return err
}

func initOptions() []grpc.ServerOption {
	options := []grpc.ServerOption{}
	unary := []grpc.UnaryServerInterceptor{
		grpc_tags.UnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	}
	stream := []grpc.StreamServerInterceptor{
		grpc_tags.StreamServerInterceptor(),
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	}

	options = append(options,
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(unary...)),
		grpc.StreamInterceptor(middleware.ChainStreamServer(stream...)),
	)

	return options
}

func (g *GRpc) Init() error {
	options := initOptions()
	g.Server = grpc.NewServer(options...)
	return nil
}

func (g *GRpc) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Config.ServerPort))
	if err != nil {
		return err
	}

	logrus.Infof("grpc server listen on %d", g.Config.ServerPort)
	if err := g.Server.Serve(lis); err != nil {
		return err
	}
	logrus.Info("grpc start success")

	return nil
}

func (g *GRpc) Close() error {
	g.Server.GracefulStop()
	return nil
}
