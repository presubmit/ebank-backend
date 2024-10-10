package api

import (
	"context"
	"net/http"
	"time"

	"ebank/shared/errors"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// RESTHandlerCallback is a callback which will be called when creating the REST server,
// and it's responsible for initializing the service implementations.
type RESTHandlerCallback func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

type Service struct {
	Address  string
	Callback RESTHandlerCallback
}

func newGateway(ctx context.Context, services []*Service) (http.Handler, error) {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(errors.ErrorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// add handlers
	dialOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			authInterceptor(),
			employeeInterceptor(),
			forwardPayloadInterceptor,
		),
		grpc.WithInsecure(),
	}
	for _, svc := range services {
		if err := svc.Callback(ctx, mux, svc.Address, dialOpts); err != nil {
			return nil, err
		}
	}

	return mux, nil
}

func applyMiddlewares(handler http.Handler) http.Handler {
	handler = extractTokenMiddleware(handler)
	handler = extractCompanyMiddleware(handler)
	handler = stripApiPrefix(handler)
	handler = handleStatusChecks(handler)
	// allowCors middleware must always be the last one
	handler = allowCORS(handler)
	return handler
}

func RunGateway(restAddress string, services []*Service) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gw, err := newGateway(ctx, services)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:         restAddress,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 30,
		Handler:      applyMiddlewares(gw),
	}

	glog.Infof("🚀 Running ApiService (%s) 🚀", restAddress)
	return srv.ListenAndServe()
}
