package main

import (
	"ebank/pb/services/auth"
	"ebank/pb/services/card"
	"ebank/pb/services/identity"
	"ebank/pb/services/ledger"
	"ebank/services/api"
	"ebank/shared/microservice"
	"flag"
	"fmt"

	"github.com/gobike/envflag"
	"github.com/golang/glog"
)

var (
	restPort = flag.String("rest_port", "80", "the port which the http server should use")
)

func main() {
	envflag.Parse()
	microservice.ParseFlags()

	services := []*api.Service{
		{
			Address:  microservice.AuthEndpoint,
			Callback: auth.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			Address:  microservice.CardEndpoint,
			Callback: card.RegisterCardServiceHandlerFromEndpoint,
		},
		{
			Address:  microservice.IdentityEndpoint,
			Callback: identity.RegisterIdentityServiceHandlerFromEndpoint,
		},
		{
			Address:  microservice.LedgerEndpoint,
			Callback: ledger.RegisterLedgerServiceHandlerFromEndpoint,
		},
	}

	address := fmt.Sprintf(":%s", *restPort)
	if err := api.RunGateway(address, services); err != nil {
		glog.Fatalf("Failed to start Api Gateway (%s) %v", address, err)
	}
}
