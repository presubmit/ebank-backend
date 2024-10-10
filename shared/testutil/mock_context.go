package testutil

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type MockContext struct {
	UserID     string
	CompanyID  string
	EmployeeID string
}

func (ctx *MockContext) Build() context.Context {
	md := make(metadata.MD)
	if ctx.UserID != "" {
		md.Set("userIdKey", ctx.UserID)
	}
	if ctx.CompanyID != "" {
		md.Set("companyIdKey", ctx.CompanyID)
	}
	if ctx.EmployeeID != "" {
		md.Set("employeeIdKey", ctx.EmployeeID)
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
