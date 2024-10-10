package main

import (
	apb "ebank/pb/services/auth"
	cpb "ebank/pb/services/card"
	ipb "ebank/pb/services/identity"
	lpb "ebank/pb/services/ledger"
	"ebank/services/api"
	"ebank/services/auth"
	"ebank/services/card"
	"ebank/services/identity"
	"ebank/services/ledger"
	"ebank/shared/microservice"
	"flag"
	"fmt"

	"github.com/gobike/envflag"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

var (
	restPort         = flag.String("rest_port", "80", "the port which the http server should use")
	jwtAccessSecret  = flag.String("jwt_access_secret", "access", "the jwt access secret key")
	jwtRefreshSecret = flag.String("jwt_refresh_secret", "refresh", "the jwt refresh secret key")
)

func main() {
	envflag.Parse()
	microservice.ParseFlags()

	// Auth
	go func() {
		microservice.RunGRPC(microservice.Options{
			ServiceName: "AuthService",
			Address:     microservice.AuthEndpoint,
			GRPCHandlerCallback: func(srv *grpc.Server) {
				apb.RegisterAuthServiceServer(srv, auth.NewService(*jwtAccessSecret, *jwtRefreshSecret))
			},
		})
	}()

	// Card
	go func() {
		microservice.RunGRPC(microservice.Options{
			ServiceName: "CardService",
			Address:     microservice.CardEndpoint,
			GRPCHandlerCallback: func(srv *grpc.Server) {
				cpb.RegisterCardServiceServer(srv, card.NewService())
			},
		})
	}()

	// Identity
	go func() {
		microservice.RunGRPC(microservice.Options{
			ServiceName: "IdentityService",
			Address:     microservice.IdentityEndpoint,
			GRPCHandlerCallback: func(srv *grpc.Server) {
				ipb.RegisterIdentityServiceServer(srv, identity.NewService())
			},
		})
	}()

	// Ledger
	go func() {
		microservice.RunGRPC(microservice.Options{
			ServiceName: "LedgerService",
			Address:     microservice.LedgerEndpoint,
			GRPCHandlerCallback: func(srv *grpc.Server) {
				lpb.RegisterLedgerServiceServer(srv, ledger.NewService())
			},
		})
	}()

	services := []*api.Service{
		{
			Address:  microservice.AuthEndpoint,
			Callback: apb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			Address:  microservice.CardEndpoint,
			Callback: cpb.RegisterCardServiceHandlerFromEndpoint,
		},
		{
			Address:  microservice.IdentityEndpoint,
			Callback: ipb.RegisterIdentityServiceHandlerFromEndpoint,
		},
		{
			Address:  microservice.LedgerEndpoint,
			Callback: lpb.RegisterLedgerServiceHandlerFromEndpoint,
		},
	}

	address := fmt.Sprintf(":%s", *restPort)
	if err := api.RunGateway(address, services); err != nil {
		glog.Fatalf("Failed to start Api Gateway (%s) %v", address, err)
	}
}
