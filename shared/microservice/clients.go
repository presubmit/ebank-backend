package microservice

import (
	apb "ebank/pb/services/auth"
	ipb "ebank/pb/services/identity"
	lpb "ebank/pb/services/ledger"

	"google.golang.org/grpc"
)

func Auth() apb.AuthServiceClient {
	conn, err := grpc.Dial(AuthEndpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return apb.NewAuthServiceClient(conn)
}

func Ledger() lpb.LedgerServiceClient {
	conn, err := grpc.Dial(LedgerEndpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return lpb.NewLedgerServiceClient(conn)
}

func Identity() ipb.IdentityServiceClient {
	conn, err := grpc.Dial(IdentityEndpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return ipb.NewIdentityServiceClient(conn)
}
