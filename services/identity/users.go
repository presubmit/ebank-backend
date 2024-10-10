package identity

import (
	"context"
	"strings"

	apb "ebank/pb/services/auth"
	pb "ebank/pb/services/identity"
	"ebank/services/identity/db"
	m "ebank/services/identity/models"
	"ebank/shared/errors"
	"ebank/shared/marqeta"
	"ebank/shared/microservice"
)

type Service struct {
	db          db.IdentityDB
	authService apb.AuthServiceClient
	marqeta     marqeta.Api
}

func NewService() *Service {
	return &Service{
		db:          db.New(),
		authService: microservice.Auth(),
		marqeta:     marqeta.NewClient(),
	}
}

func (s *Service) RegisterUser(ctx context.Context, r *pb.RegisterUserRequest) (*pb.AuthResponse, error) {
	u := &m.User{
		Email:     r.GetEmail(),
		Password:  r.GetPassword(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
	}
	if err := u.ValidateFields(); err != nil {
		return nil, errors.InvalidArgument(err)
	}
	if err := u.HashPassword(); err != nil {
		return nil, errors.Internal(err)
	}

	// save user in db
	userId, err := s.db.CreateUser(nil, u)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, errors.InvalidArgumentf("email already registered")
		}
		return nil, errors.Internal(err)
	}

	// generate access tokens using AuthService
	tokenRes, err := s.authService.GenerateToken(ctx, &apb.GenerateTokenRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Internal(err)
	}
	return &pb.AuthResponse{
		AccessToken:  tokenRes.GetAccessToken(),
		RefreshToken: tokenRes.GetRefreshToken(),
	}, nil
}

func (s *Service) LoginUser(ctx context.Context, r *pb.LoginUserRequest) (*pb.AuthResponse, error) {
	u, err := s.db.GetUserByEmail(nil, r.GetEmail())
	if err != nil {
		return nil, errors.NotAuthenticatedf("invalid credentials")
	}

	if err := u.ComparePassword(r.GetPassword()); err != nil {
		return nil, errors.NotAuthenticatedf("invalid credentials")
	}

	// generate access tokens using AuthService
	tokenRes, err := s.authService.GenerateToken(ctx, &apb.GenerateTokenRequest{
		UserId: u.ID,
	})
	if err != nil {
		return nil, errors.Internal(err)
	}
	return &pb.AuthResponse{
		AccessToken:  tokenRes.GetAccessToken(),
		RefreshToken: tokenRes.GetRefreshToken(),
	}, nil
}

func (s *Service) GetCurrentUser(ctx context.Context, r *pb.Empty) (*pb.User, error) {
	userId, err := microservice.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	u, err := s.db.GetUser(nil, userId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	return u.ToProto(), nil
}
