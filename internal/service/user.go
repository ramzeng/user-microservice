package service

import (
	"context"
	"github.com/ramzeng/user-microservice/internal/model"
	"github.com/ramzeng/user-microservice/pb"
	"github.com/ramzeng/user-microservice/pkg/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	pb.UnimplementedUserServer
}

func (u *User) CreateUser(context context.Context, request *pb.CreateUserRequest) (*pb.UserResponse, error) {
	var user model.User

	database.Where("email = ?", request.Email).First(&user)

	if user.ID != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "email %s already exists", request.Email)
	}

	user.Email = request.Email
	user.Password = request.Password

	database.Create(&user)

	return &pb.UserResponse{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}, nil
}
