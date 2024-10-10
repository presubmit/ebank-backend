package microservice

import (
	"fmt"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCHandlerCallback is a callback which will be called when creating the GRPC server,
// and it's responsible for initializing the service implementations.
type GRPCHandlerCallback func(*grpc.Server)

type Options struct {
	ServiceName         string
	GRPCHandlerCallback GRPCHandlerCallback
	Address             string
}

func startGRPCServer(opts Options) error {
	lis, err := net.Listen("tcp", opts.Address)
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	// call gRPC handler callback
	opts.GRPCHandlerCallback(srv)
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)

	return srv.Serve(lis)
}

func RunGRPC(opts Options) {
	if len(opts.Address) == 0 {
		opts.Address = fmt.Sprintf(":%s", *grpcPort)
	}

	glog.Infof("🚀 Running %s GRPC(%s) 🚀", opts.ServiceName, opts.Address)
	if err := startGRPCServer(opts); err != nil {
		glog.Fatalf("Failed to start GRPC %s : %s", opts.ServiceName, err)
	}
}
