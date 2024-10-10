package main

import (
	pb "ebank/pb/services/ledger"
	"ebank/services/ledger"
	"ebank/shared/microservice"

	"google.golang.org/grpc"
)

func main() {
	microservice.ParseFlags()
	s := ledger.NewService()
	// Start GRPC service
	microservice.RunGRPC(microservice.Options{
		ServiceName: "LedgerService",
		GRPCHandlerCallback: func(srv *grpc.Server) {
			pb.RegisterLedgerServiceServer(srv, s)
		},
	})
}
