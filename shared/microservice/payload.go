package microservice

import (
	"context"
	"ebank/shared/errors"

	"google.golang.org/grpc/metadata"
)

const (
	userIdKey     = "userIdKey"
	companyIdKey  = "companyIdKey"
	employeeIdKey = "employeeIdKey"
)

func GetUserId(ctx context.Context) (string, error) {
	return mdValue(ctx, userIdKey)
}

func GetCompanyId(ctx context.Context) (string, error) {
	return mdValue(ctx, companyIdKey)
}

func GetEmployeeId(ctx context.Context) (string, error) {
	return mdValue(ctx, employeeIdKey)
}

func mdValue(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.Internalf("couldn't get incoming context")
	}
	vals := md.Get(key)
	if len(vals) == 0 {
		return "", errors.Internalf("metadata doesn't contain key: %s", key)
	}
	if len(vals) > 1 {
		return "", errors.Internalf("metadata contains multiple values for key: %s", key)
	}
	return vals[0], nil
}
