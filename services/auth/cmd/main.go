package main

import (
	"flag"

	pb "ebank/pb/services/auth"
	"ebank/services/auth"
	"ebank/shared/microservice"

	"github.com/gobike/envflag"
	"google.golang.org/grpc"
)

var (
	jwtAccessSecret  = flag.String("jwt_access_secret", "access", "the jwt access secret key")
	jwtRefreshSecret = flag.String("jwt_refresh_secret", "refresh", "the jwt refresh secret key")
)

func main() {
	envflag.Parse()
	microservice.ParseFlags()

	s := auth.NewService(*jwtAccessSecret, *jwtRefreshSecret)
	// Start GRPC service
	microservice.RunGRPC(microservice.Options{
		ServiceName: "AuthService",
		GRPCHandlerCallback: func(srv *grpc.Server) {
			pb.RegisterAuthServiceServer(srv, s)
		},
	})
}
