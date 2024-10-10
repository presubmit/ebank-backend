package microservice

import (
	"flag"

	"github.com/gobike/envflag"
)

var (
	grpcPort = flag.String("grpc_port", "8079", "the port which the grpc server should use")

	AuthEndpoint     string
	CardEndpoint     string
	IdentityEndpoint string
	LedgerEndpoint   string

	DbHost string
	DbPort string
	DbPass string
	DbUser string
	DbName string

	RedisHost string
	RedisPort string
)

func ParseFlags() {
	_ = flag.Set("alsologtostderr", "true")

	flag.StringVar(&AuthEndpoint, "auth_service", "authsvc:8079", "the address of AuthService")
	flag.StringVar(&CardEndpoint, "card_service", "cardsvc:8079", "the address of CardService")
	flag.StringVar(&IdentityEndpoint, "identity_service", "identitysvc:8079", "the address of IdentityService")
	flag.StringVar(&LedgerEndpoint, "ledger_service", "ledgersvc:8079", "the address of LedgerService")

	flag.StringVar(&DbHost, "db_host", "localhost", "database host")
	flag.StringVar(&DbPort, "db_port", "5432", "database port")
	flag.StringVar(&DbPass, "db_pass", "ebank123", "database password")
	flag.StringVar(&DbUser, "db_user", "postgres", "database username")
	flag.StringVar(&DbName, "db_name", "ebank", "database name")

	flag.StringVar(&RedisHost, "redis_host", "redis", "redis host address")
	flag.StringVar(&RedisPort, "redis_port", "6379", "redis port address")

	envflag.Parse()
}
