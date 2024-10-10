package identity

import (
	"context"
	apb "ebank/pb/services/auth"
	pb "ebank/pb/services/identity"
	"ebank/services/identity/db"
	m "ebank/services/identity/models"
	"ebank/shared/errors"
	"testing"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		req                *pb.RegisterUserRequest
		userID             string
		createUserErr      error
		tokenRes           *apb.TokenResponse
		generateTokenErr   error
		generateTokenCalls int32
		res                *pb.AuthResponse
		err                error
	}{
		{
			req: &pb.RegisterUserRequest{
				Email:     "johndoe@test.com",
				Password:  "123456",
				FirstName: "John",
				LastName:  "Doe",
			},
			userID:             "test-1",
			tokenRes:           &apb.TokenResponse{AccessToken: "123", RefreshToken: "456"},
			generateTokenCalls: 1,
			res:                &pb.AuthResponse{AccessToken: "123", RefreshToken: "456"},
		},
		{
			req: &pb.RegisterUserRequest{
				Email:     "johndoe@test.com",
				Password:  "123456",
				FirstName: "John",
				LastName:  "Doe",
			},
			userID:             "test-1",
			generateTokenErr:   errors.Internalf("token error"),
			generateTokenCalls: 1,
			err:                errors.Internalf("token error"),
		},
		{
			req: &pb.RegisterUserRequest{
				Email:     "johndoe@test.com",
				Password:  "123456",
				FirstName: "John",
				LastName:  "Doe",
			},
			createUserErr: errors.Internalf("db error"),
			err:           errors.Internalf("db error"),
		},
		{
			req: &pb.RegisterUserRequest{
				Password:  "123456",
				FirstName: "John",
				LastName:  "Doe",
			},
			err: errors.InvalidArgument(),
		},
	}

	for _, test := range tests {
		mockAuthService := &apb.MockAuthServiceClient{
			GenerateTokenFunc: func(r *apb.GenerateTokenRequest) (*apb.TokenResponse, error) {
				if r.GetUserId() != test.userID {
					t.Errorf("Register(%v): authservice.GenerateToken expected user id %s; got %s", test.req, test.userID, r.GetUserId())
				}
				return test.tokenRes, test.generateTokenErr
			},
		}
		s := &Service{
			db: &db.MockIdentityDB{
				CreateUserFunc: func(*m.User) (string, error) {
					return test.userID, test.createUserErr
				},
			},
			authService: mockAuthService,
		}
		got, err := s.RegisterUser(context.Background(), test.req)
		if mockAuthService.GenerateTokenCalls != test.generateTokenCalls {
			t.Errorf("Register(%v) didn't call authservice.GenerateToken the right number of times; expected %v; actual %v", test.req, test.generateTokenCalls, mockAuthService.GenerateTokenCalls)
		}
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("Register(%v) didn't return error; expected %v; got %v", test.req, status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("Register(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		req                *pb.LoginUserRequest
		user               *m.User
		getUserErr         error
		tokenRes           *apb.TokenResponse
		generateTokenErr   error
		generateTokenCalls int32
		res                *pb.AuthResponse
		err                error
	}{
		{
			req: &pb.LoginUserRequest{
				Email:    "johndoe@test.com",
				Password: "testpass",
			},
			user:               &m.User{ID: "id-1", Password: "testpass"},
			tokenRes:           &apb.TokenResponse{AccessToken: "123", RefreshToken: "456"},
			generateTokenCalls: 1,
			res:                &pb.AuthResponse{AccessToken: "123", RefreshToken: "456"},
		},
		{
			req: &pb.LoginUserRequest{
				Email:    "johndoe@test.com",
				Password: "testpass",
			},
			user:               &m.User{ID: "id-1", Password: "otherpass"},
			generateTokenCalls: 0,
			err:                errors.NotAuthenticated(),
		},
		{
			req: &pb.LoginUserRequest{
				Email:    "johndoe@test.com",
				Password: "testpass",
			},
			getUserErr:         errors.Internal(),
			generateTokenCalls: 0,
			err:                errors.NotAuthenticated(),
		},
		{
			req: &pb.LoginUserRequest{
				Email:    "johndoe@test.com",
				Password: "testpass",
			},
			user:               &m.User{ID: "id-1", Password: "testpass"},
			generateTokenErr:   errors.Internalf("token error"),
			generateTokenCalls: 1,
			err:                errors.Internal(),
		},
	}

	for _, test := range tests {
		mockAuthService := &apb.MockAuthServiceClient{
			GenerateTokenFunc: func(r *apb.GenerateTokenRequest) (*apb.TokenResponse, error) {
				if r.GetUserId() != test.user.ID {
					t.Errorf("Register(%v): authservice.GenerateToken expected user id %s; got %s", test.req, test.user.ID, r.GetUserId())
				}
				return test.tokenRes, test.generateTokenErr
			},
		}
		s := &Service{
			db: &db.MockIdentityDB{
				GetUserByEmailFunc: func(_ string) (*m.User, error) {
					if test.user != nil {
						_ = test.user.HashPassword()
					}
					return test.user, test.getUserErr
				},
			},
			authService: mockAuthService,
		}
		got, err := s.LoginUser(context.Background(), test.req)
		if mockAuthService.GenerateTokenCalls != test.generateTokenCalls {
			t.Errorf("Login(%v) didn't call authservice.GenerateToken the right number of times; expected %v; actual %v", test.req, test.generateTokenCalls, mockAuthService.GenerateTokenCalls)
		}
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("Login(%v) didn't return error; expected %v; got %v", test.req, status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("Login(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}
