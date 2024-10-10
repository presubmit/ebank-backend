package main

import (
	pb "ebank/pb/services/identity"
	"ebank/services/identity"
	"ebank/shared/microservice"

	"google.golang.org/grpc"
)

func main() {
	microservice.ParseFlags()
	s := identity.NewService()
	// Start GRPC service
	microservice.RunGRPC(microservice.Options{
		ServiceName: "IdentityService",
		GRPCHandlerCallback: func(srv *grpc.Server) {
			pb.RegisterIdentityServiceServer(srv, s)
		},
	})
}
