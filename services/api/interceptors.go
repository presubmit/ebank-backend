package api

import (
	"context"
	apb "ebank/pb/services/auth"
	ipb "ebank/pb/services/identity"
	"ebank/shared/errors"
	"ebank/shared/microservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

var (
	userIdKeyStr       = "userIdKey"
	userIdKey          = contextKey(userIdKeyStr)
	companyIdKeyStr    = "companyIdKey"
	companyIdKey       = contextKey(companyIdKeyStr)
	employeeIdKeyStr   = "employeeIdKey"
	employeeIdKey      = contextKey(employeeIdKeyStr)
	employeeRoleKeyStr = "employeeRoleKey"
	employeeRoleKey    = contextKey(employeeRoleKeyStr)
	tokenKey           = contextKey("tokenKey")
)

func forwardPayloadInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.MD{}
	if userId := ctx.Value(userIdKey); userId != nil {
		md.Set(userIdKeyStr, userId.(string))
	}
	if companyId := ctx.Value(companyIdKey); companyId != nil {
		md.Set(companyIdKeyStr, companyId.(string))
	}
	if employeeId := ctx.Value(employeeIdKey); employeeId != nil {
		md.Set(employeeIdKeyStr, employeeId.(string))
	}
	if employeeRole := ctx.Value(employeeRoleKey); employeeRole != nil {
		md.Set(employeeRoleKeyStr, employeeRole.(string))
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}

var skipAuth = map[string]bool{
	"/auth.AuthService/RefreshToken":         true,
	"/identity.IdentityService/LoginUser":    true,
	"/identity.IdentityService/RegisterUser": true,
}

func authInterceptor() grpc.UnaryClientInterceptor {
	authService := microservice.Auth()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if _, ok := skipAuth[method]; ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		token := ctx.Value(tokenKey)
		if token == nil {
			return errors.NotAuthenticatedf("not authorized")
		}

		tokenRes, err := authService.VerifyToken(ctx, &apb.VerifyTokenRequest{
			AccessToken: token.(string),
		})
		if err != nil {
			return errors.NotAuthenticatedf("not authorized")
		}
		ctx = context.WithValue(ctx, userIdKey, tokenRes.UserId)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

var skipEmployeeCheck = map[string]bool{
	"/identity.IdentityService/CreateCompany":  true,
	"/identity.IdentityService/GetCompanies":   true,
	"/identity.IdentityService/GetCurrentUser": true,
}

func employeeInterceptor() grpc.UnaryClientInterceptor {
	identityService := microservice.Identity()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if _, ok := skipEmployeeCheck[method]; ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		userId := ctx.Value(userIdKey)
		if userId == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		companyId := ctx.Value(companyIdKey)
		if companyId == nil {
			return errors.InvalidArgumentf("companyId missing")
		}

		emp, err := identityService.VerifyEmployee(ctx, &ipb.Employee{
			CompanyId: companyId.(string),
			UserId:    userId.(string),
		})
		if err != nil {
			return errors.InvalidArgumentf("invalid company")
		}
		ctx = context.WithValue(ctx, employeeIdKey, emp.Id)
		ctx = context.WithValue(ctx, employeeRoleKey, emp.Role)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
