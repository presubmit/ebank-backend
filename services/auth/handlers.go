package auth

import (
	"context"
	pb "ebank/pb/services/auth"
	"ebank/services/auth/db"
	"ebank/services/auth/tokens"
	er "ebank/shared/errors"
	"ebank/shared/redis"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	AccessTokenValidity  = 20 * time.Hour // 20 hours
	RefreshTokenValidity = 2 * time.Hour  // 2 hours
)

type Service struct {
	db               db.AuthDB
	cache            redis.Redis
	jwtAccessSecret  []byte
	jwtRefreshSecret []byte
}

func NewService(jwtAccessSecret, jwtRefreshSecret string) *Service {
	return &Service{
		db:               db.New(),
		cache:            redis.New(),
		jwtAccessSecret:  []byte(jwtAccessSecret),
		jwtRefreshSecret: []byte(jwtRefreshSecret),
	}
}

func (s *Service) VerifyToken(ctx context.Context, r *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	td, err := tokens.ParseAccessToken(r.GetAccessToken(), s.jwtAccessSecret)
	if err != nil {
		return nil, err
	}
	// extract user id from cache for the given token uuid
	userId, err := s.cache.Get(td.ID).Result()
	if err != nil {
		return nil, err
	}
	// successfully authenticated
	return &pb.VerifyTokenResponse{UserId: userId}, nil
}

func (s *Service) GenerateToken(ctx context.Context, r *pb.GenerateTokenRequest) (*pb.TokenResponse, error) {
	userId := r.GetUserId()
	if len(userId) == 0 {
		return nil, errors.New("invalid user id")
	}
	// define access token
	at := &tokens.AccessToken{
		ID:        uuid.New().String(),
		UserID:    userId,
		ExpiresAt: time.Now().Add(AccessTokenValidity),
	}
	// generate jwt access token
	accessToken, err := at.GenerateJWT(s.jwtAccessSecret)
	if err != nil {
		return nil, err
	}
	// save access token in cache
	if err := s.cache.Set(at.ID, at.UserID, AccessTokenValidity).Err(); err != nil {
		return nil, err
	}

	// define refresh token
	rt := &tokens.RefreshToken{
		ID:            uuid.New().String(),
		UserID:        userId,
		AccessTokenID: at.ID,
		ExpiresAt:     time.Now().Add(RefreshTokenValidity),
	}
	// generate jwt refresh token
	refreshToken, err := rt.GenerateJWT(s.jwtRefreshSecret)
	if err != nil {
		return nil, err
	}
	// save refresh token in db
	if err := s.db.SaveRefreshToken(nil, rt); err != nil {
		return nil, err
	}

	return &pb.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, r *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	tkn, err := tokens.ParseRefreshToken(r.GetRefreshToken(), s.jwtRefreshSecret)
	if err != nil {
		return nil, er.NotAuthenticatedf("invalid refresh token")
	}
	// if there is an active access token, we shouldn't exchange the refresh token
	if _, err := s.cache.Get(tkn.AccessTokenID).Result(); err == nil {
		return nil, er.NotAuthenticatedf("access token is still valid")
	}

	// check if refresh token is in db, if it's valid and if the access uuid matches the one in the jwt
	rt, err := s.db.GetRefreshToken(nil, tkn.ID)
if err != nil || rt.AccessTokenID != tkn.AccessTokenID || rt.ExpiresAt.Before(time.Now()) {
		return nil, er.NotAuthenticatedf("invalid refresh token")
	}

	// delete existing refresh token
	if err := s.db.DeleteRefreshToken(nil, rt.ID); err != nil {
		return nil, er.Internal(err)
	}

	// generate a new set of access and refresh tokens
	return s.GenerateToken(ctx, &pb.GenerateTokenRequest{
		UserId: rt.UserID,
	})
}
