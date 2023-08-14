package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ramzeng/user-microservice/internal/core"
	"github.com/ramzeng/user-microservice/internal/model"
	"github.com/ramzeng/user-microservice/internal/service"
	"github.com/ramzeng/user-microservice/pb"
	"github.com/ramzeng/user-microservice/pkg/config"
	"github.com/ramzeng/user-microservice/pkg/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
)

func main() {
	configPath := flag.String("path", "", "the path of config file")
	migrate := flag.Bool("migrate", false, "auto migrate database")

	flag.Parse()

	core.ReadConfig(*configPath)

	core.Initialize()

	if *migrate {
		if err := database.AutoMigrate(&model.User{}); err != nil {
			panic(fmt.Sprintf("auto migrate database failed: %s", err))
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetInt("service.port")))

	if err != nil {
		panic(fmt.Sprintf("tcp listen failed: %s", err))
	}

	server := grpc.NewServer()

	pb.RegisterAuthServer(server, &service.Auth{})
	pb.RegisterUserServer(server, &service.User{})

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(fmt.Sprintf("grpc server start failed: %s", err))
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = pb.RegisterAuthHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.GetInt("service.port")), opts)

	if err != nil {
		panic(fmt.Sprintf("register auth handler from endpoint failed: %s", err))
	}

	err = pb.RegisterUserHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.GetInt("service.port")), opts)

	if err != nil {
		panic(fmt.Sprintf("register user handler from endpoint failed: %s", err))
	}

	err = http.ListenAndServe(":9090", mux)

	if err != nil {
		panic(fmt.Sprintf("http listen and serve failed: %s", err))
	}
}
