package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorBody struct {
	Error string `json:"error,omitempty"`
}

func ErrorHandler(ctx context.Context, _ *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, _ *http.Request, e error) {
	w.Header().Set("Content-type", m.ContentType(w))

	var msg string
	switch status.Code(e) {
	case codes.NotFound:
		msg = "not found"
	case codes.Internal:
		// don't expose internal errors
		msg = "internal error"
	default:
		msg = status.Convert(e).Message()
	}

	if status.Code(e) == codes.Internal {
		glog.Errorf(fmt.Sprintf("%v: %s", codes.Internal, status.Convert(e).Message()))
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(e)))
	_ = json.NewEncoder(w).Encode(errorBody{msg})
}

func newError(c codes.Code, str string, args ...interface{}) error {
	return status.Errorf(c, str, args...)
}

func NotAuthenticated(err ...error) error {
	e := ""
	if len(err) > 0 {
		e = err[0].Error()
	}
	return NotAuthenticatedf(e)
}

func NotAuthenticatedf(str string, args ...interface{}) error {
	return newError(codes.Unauthenticated, str, args...)
}

func Internal(err ...error) error {
	e := ""
	if len(err) > 0 {
		e = err[0].Error()
	}
	return Internalf(e)
}

func Internalf(str string, args ...interface{}) error {
	return newError(codes.Internal, str, args...)
}

func Duplicate(err ...error) error {
	e := ""
	if len(err) > 0 {
		e = err[0].Error()
	}
	return Duplicatef(e)
}

func Duplicatef(str string, args ...interface{}) error {
	return newError(codes.AlreadyExists, str, args...)
}

func NotFound(err ...error) error {
	e := ""
	if len(err) > 0 {
		e = err[0].Error()
	}
	return NotFoundf(e)
}

func NotFoundf(str string, args ...interface{}) error {
	return newError(codes.NotFound, str, args...)
}

func InvalidArgument(err ...error) error {
	e := ""
	if len(err) > 0 {
		e = err[0].Error()
	}
	return InvalidArgumentf(e)
}

func InvalidArgumentf(str string, args ...interface{}) error {
	return newError(codes.InvalidArgument, str, args...)
}
