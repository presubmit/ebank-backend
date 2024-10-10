package main

import (
	pb "ebank/pb/services/card"
	"ebank/services/card"
	"ebank/shared/microservice"

	"google.golang.org/grpc"
)

func main() {
	microservice.ParseFlags()
	s := card.NewService()
	// Start GRPC service
	microservice.RunGRPC(microservice.Options{
		ServiceName: "CardService",
		GRPCHandlerCallback: func(srv *grpc.Server) {
			pb.RegisterCardServiceServer(srv, s)
		},
	})
}
