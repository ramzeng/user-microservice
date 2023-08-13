package main

import (
	"context"
	"fmt"
	"github.com/ramzeng/user-microservice/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(fmt.Sprintf("grpc connection failed: %s", err))
	}

	defer func() {
		_ = connection.Close()
	}()

	userClient := pb.NewUserClient(connection)

	userResponse, err := userClient.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email:    "email@example.com",
		Password: "password",
	})

	if err != nil {
		panic(fmt.Sprintf("create user failed: %s", err))
	}

	fmt.Println("====== create user successfully ======")
	fmt.Printf("user ID: %d\n", userResponse.Id)
	fmt.Printf("user Email: %s\n", userResponse.Email)
	fmt.Printf("user CreatedAt: %d\n", userResponse.CreatedAt)
	fmt.Printf("user UpdatedAt: %d\n", userResponse.UpdatedAt)
	fmt.Println()

	authClient := pb.NewAuthClient(connection)

	tokenResponse, err := authClient.CreateTokenViaPassword(context.Background(), &pb.CreateTokenViaPasswordRequest{
		Email:    "email@example.com",
		Password: "password",
	})

	if err != nil {
		panic(fmt.Sprintf("create token failed: %s", err))
	}

	fmt.Println("====== create token successfully ======")
	fmt.Printf("Token: %s\n", tokenResponse.AccessToken)
	fmt.Printf("Token ExpiredAt: %d\n", tokenResponse.AccessTokenExpiredAt)
	fmt.Printf("Refresh Token: %s\n", tokenResponse.RefreshToken)
	fmt.Printf("Refresh Token ExpiredAt: %d\n", tokenResponse.RefreshTokenExpiredAt)
	fmt.Println()

	userResponse, err = authClient.GetUserViaAccessToken(context.Background(), &pb.GetUserViaAccessTokenRequest{
		AccessToken: tokenResponse.AccessToken,
	})

	if err != nil {
		panic(fmt.Sprintf("get user failed: %s", err))
	}

	fmt.Println("====== get user successfully ======")
	fmt.Printf("user ID: %d\n", userResponse.Id)
	fmt.Printf("user Email: %s\n", userResponse.Email)
	fmt.Printf("user CreatedAt: %d\n", userResponse.CreatedAt)
	fmt.Printf("user UpdatedAt: %d\n", userResponse.UpdatedAt)
	fmt.Println()
}
