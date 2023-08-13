package service

import (
	"context"
	"github.com/ramzeng/user-microservice/internal/model"
	"github.com/ramzeng/user-microservice/pb"
	"github.com/ramzeng/user-microservice/pkg/auth"
	"github.com/ramzeng/user-microservice/pkg/database"
	"github.com/ramzeng/user-microservice/pkg/logger"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct {
	pb.UnimplementedAuthServer
}

func (a *Auth) GetUserViaAccessToken(ctx context.Context, request *pb.GetUserViaAccessTokenRequest) (*pb.UserResponse, error) {
	if !auth.VerifyToken(request.AccessToken) {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	subject, err := auth.ParseTokenSubject(request.AccessToken)

	if err != nil {
		logger.Info("app", "[auth] | token parse failed: %s", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	var user model.User

	database.Where("id = ?", subject).First(&user)

	return &pb.UserResponse{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}, nil
}

func (a *Auth) CreateTokenViaPassword(ctx context.Context, request *pb.CreateTokenViaPasswordRequest) (*pb.TokenResponse, error) {
	var user model.User

	database.Where("email = ?", request.Email).First(&user)

	if user.ID == 0 {
		return nil, status.Error(codes.PermissionDenied, "email or password incorrect")
	}

	// TODO: use bcrypt to compare password
	if user.Password != request.Password {
		return nil, status.Error(codes.PermissionDenied, "email or password incorrect")
	}

	token, err := auth.CreateToken(cast.ToString(user.ID), true)

	if err != nil {
		logger.Error("app", "[auth] | create token failed: %s", zap.Error(err))
		return nil, status.Error(codes.Internal, "create token failed")
	}

	return &pb.TokenResponse{
		AccessToken:           token.AccessToken,
		AccessTokenExpiredAt:  token.AccessTokenExpiredAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiredAt: token.RefreshTokenExpiredAt,
	}, nil
}

func (a *Auth) CreateTokenViaRefreshToken(ctx context.Context, request *pb.CreateTokenViaRefreshTokenRequest) (*pb.TokenResponse, error) {
	if !auth.VerifyToken(request.RefreshToken) {
		return nil, status.Error(codes.PermissionDenied, "invalid refresh token")
	}

	if !auth.MatchAccessTokenAndRefreshToken(request.AccessToken, request.RefreshToken) {
		return nil, status.Error(codes.PermissionDenied, "invalid refresh token")
	}

	token, err := auth.RefreshAccessToken(request.AccessToken)

	if err != nil {
		logger.Error("app", "[auth] | refresh token failed: %s", zap.Error(err))
		return nil, status.Error(codes.Internal, "refresh token failed")
	}

	return &pb.TokenResponse{
		AccessToken:          token.AccessToken,
		AccessTokenExpiredAt: token.AccessTokenExpiredAt,
	}, nil
}
