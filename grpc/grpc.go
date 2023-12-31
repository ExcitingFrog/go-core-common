package grpc

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/ExcitingFrog/go-core-common/log"
	"github.com/ExcitingFrog/go-core-common/provider"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRpc struct {
	provider.IProvider

	addr    string
	Server  *grpc.Server
	Config  *Config
	Running bool
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
	log.Logger().With(
		zap.String("stack", string(debug.Stack())),
		zap.String("error", err.Error()),
	).Error("grpc service panic")
	return err
}

func initOptions() []grpc.ServerOption {
	options := []grpc.ServerOption{}
	unary := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
		grpc_tags.UnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	}
	stream := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
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
	g.addr = fmt.Sprintf(":%d", g.Config.ServerPort)
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}

	reflection.Register(g.Server)

	log.Logger().With(
		zap.String("port", g.addr),
	).Info("grpc server start")

	g.Running = true
	if err := g.Server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (g *GRpc) Close() error {
	g.Server.GracefulStop()
	return nil
}
